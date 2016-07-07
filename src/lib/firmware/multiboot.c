/*-
 * Copyright (c) 2016 Thomas Haggett
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THOMAS HAGGETT ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL NETAPP, INC OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <xhyve/firmware/multiboot.h>
#include <xhyve/vmm/vmm_api.h>

#define MULTIBOOT_MAGIC 0x1BADB002
#define MULTIBOOT_SEARCH_END 0x2000

struct multiboot_header {
  uint32_t magic;
  uint32_t flags;
  uint32_t checksum;
};

struct multiboot_load_header {
  uint32_t header_addr;
  uint32_t load_addr;
  uint32_t load_end_addr;
  uint32_t bss_end_addr;
  uint32_t entry_addr;
};

struct multiboot_info  {
  uint32_t flags;
  uint32_t mem_lower;
  uint32_t mem_upper;
  uint32_t boot_device;
  uint32_t cmdline_addr;
  uint32_t mods_count;
  uint32_t mods_addr;
};

struct multiboot_module_entry {
  uint32_t addr_start;
  uint32_t addr_end;
  uint32_t cmdline;
  uint32_t pad;
};

static struct multiboot_config {
  char* kernel_path;
  char* module_list;
  char* kernel_append;
} config;

struct boot_config {
  long header_offset;
  uint16_t load_alignment;
  char provide_mem_headers;
  char padding1[5];
  uintptr_t guest_mem_base;
  uintptr_t guest_mem_size;

  char padding2[4];
  struct multiboot_load_header kernel_load_data;
};

// #define ALIGN(p,n) (void*)((((uintptr_t)p/n) + 1)*n)

int multiboot_find_header(FILE* image, struct boot_config* boot_config);
size_t multiboot_load_image(FILE* image, struct boot_config *boot_config);
uint64_t multiboot_set_guest_state(struct boot_config* boot_config);

//
// called by xhyve to pass in the firmware arguments
//
void multiboot_init(char *kernel_path, char *module_list, char *kernel_append) {
  config.kernel_path = kernel_path;
  config.module_list = module_list;
  config.kernel_append = kernel_append;
}

// 
// scans the configured kernel for it's multiboot header.
// returns 0 if no multiboot header is found, non-zero if one is.
int multiboot_find_header(FILE* image, struct boot_config* boot_config) {
  struct multiboot_header header;
  uint8_t found = 0;

  fseek(image, 0L, SEEK_SET);

  while(!found || !feof(image) || ftell(image) < MULTIBOOT_SEARCH_END) {

    boot_config->header_offset = ftell(image);
    // fill our header struct with data from the file
    if( 1 != fread(&header, sizeof(struct multiboot_header), 1, image)) {
      perror("error reading from kernel");
      return 0;
    }
    
    if(header.magic == MULTIBOOT_MAGIC && header.checksum + header.flags + header.magic == 0) {
      found = 1;
      break;
    }
      
    // make sure to jump 64-bits back in the file, so that if there happens to be two
    // magic values one after another, the second one doesn't get missed. The only 
    // requirement is the header is 32-bit aligned.
    fseek(image, -8, SEEK_CUR);
  }

  if( found == 0 ) return 0;

  printf("Parsing multiboot header:\n");

  // are there any mandatory flags that we don't support? (any other than 0 and 1 set)
  uint16_t supported_mandatory = ((1<<1) | (1<<0));
  if( ((header.flags & ~supported_mandatory) & 0xFFFF) != 0x0) {
    printf("Multiboot header has unsupported mandatory flags, bailing.\n");
    return 0;
  }

  // at this point, we need to check the flags and pull in the additional sections
  if( header.flags & (1<<0) ) {
    printf(" (bit 0) Loading modules with 4K alignment\n");
    boot_config->load_alignment = 4096;
  } else {
    boot_config->load_alignment = 1;
  }

  if( header.flags & (1<<1) ) {
    printf(" (bit 1) Must provide mem_* fields in multiboot header\n");
    boot_config->provide_mem_headers = 1;
  } else {
    boot_config->provide_mem_headers = 0;
  }

  if( header.flags & (1<<16) ) {
    printf(" (bit 16) Multiboot image placement fields are valid\n");

    // read the memory placement header directly into the boot config struct
    if( 1 != fread((void*)&boot_config->kernel_load_data, sizeof(struct multiboot_load_header), 1, image) ) {
      perror("failed to read image placement data from multiboot header");
      return 0;
    }

  } else {
    printf(" (no bit 16) Image placement fields aren't valid - TODO: still try to load?\n");
    return 0;
  }
  
  return found;
}

void* guest_to_host(void* guest_addr, struct boot_config *boot_config);
void* host_to_guest(void* host_addr, struct boot_config *boot_config);
void* guest_to_host(void* guest_addr, struct boot_config *boot_config) {
  return (void*)(boot_config->guest_mem_base + (uintptr_t)guest_addr);
}
void* host_to_guest(void* host_addr, struct boot_config *boot_config) {
  return (void*)((uintptr_t)host_addr - boot_config->guest_mem_base);
}

size_t multiboot_load_image(FILE* image, struct boot_config *boot_config) {

  size_t image_load_size;
  unsigned long load_offset = (unsigned long)(boot_config->kernel_load_data.header_addr - boot_config->header_offset);

  // if there wasn't a load_end_addr provided, then default it to the length of the image file
  if( boot_config->kernel_load_data.load_end_addr == 0x0 ) {
    fseek(image, 0x0, SEEK_END);
    boot_config->kernel_load_data.load_end_addr = (uint32_t) ftell(image) + (uint32_t)load_offset;
  }
  image_load_size = boot_config->kernel_load_data.load_end_addr - boot_config->kernel_load_data.load_addr;
  printf("image load size is %zu\n", image_load_size);
  printf("header addr   = %x\n", boot_config->kernel_load_data.header_addr);
  printf("load addr     = %x\n", boot_config->kernel_load_data.load_addr);
  printf("load end addr = %x\n", boot_config->kernel_load_data.load_end_addr);
  printf("bss end addr  = %x\n", boot_config->kernel_load_data.bss_end_addr);
  printf("entry addr    = %x\n", boot_config->kernel_load_data.entry_addr);

  printf("Header offset is %lu, wants physical address %u\n", boot_config->header_offset, boot_config->kernel_load_data.header_addr);

  // Jump to the load offset
  printf("Load addr is %i, offset is %lu seeking to %li in file\n", boot_config->kernel_load_data.load_addr,load_offset, boot_config->kernel_load_data.load_addr - (long)load_offset);
  if(0 != fseek(image, boot_config->kernel_load_data.load_addr - (long)load_offset, SEEK_SET)) {
    perror("fseek() to specified load_addr failed on the kernel image");
    return 0;
  }
  
  if( 1 != fread(guest_to_host((void *)(uintptr_t)boot_config->kernel_load_data.load_addr, boot_config), image_load_size, 1, image)) {
    perror("fread() kernel image");
    return 0;
  }

  void* magic_addr = guest_to_host((void*)(uintptr_t)boot_config->kernel_load_data.header_addr, boot_config);
  printf("sanity-check: this should be the magic value: %x\n", *(uint32_t*)magic_addr);

  return image_load_size;
}

uint64_t multiboot_set_guest_state(struct boot_config* boot_config) {

  struct multiboot_header* header = (struct multiboot_header*)boot_config->guest_mem_base;
  void* guest_header_ptr = host_to_guest(header, boot_config);
  header->flags = 0x1234;


  xh_vcpu_reset(0);
  xh_vm_set_register(0, VM_REG_GUEST_CR0, 0x21);
  xh_vm_set_register(0, VM_REG_GUEST_RAX, 0x2BADB002);
  xh_vm_set_register(0, VM_REG_GUEST_RBX, (uint64_t)guest_header_ptr);
  xh_vm_set_register(0, VM_REG_GUEST_RIP, boot_config->kernel_load_data.entry_addr);

  // xh_vm_set_desc(0, VM_REG_GUEST_GDTR, (uintptr_t)gdt_entry - (uintptr_t)gpa_map, 0x1f, 0);
  // xh_vm_set_desc(0, VM_REG_GUEST_CS, 0, 0xffffffff, 0xc09b);
  // xh_vm_set_desc(0, VM_REG_GUEST_DS, 0, 0xffffffff, 0xc093);
  // xh_vm_set_desc(0, VM_REG_GUEST_ES, 0, 0xffffffff, 0xc093);
  // xh_vm_set_desc(0, VM_REG_GUEST_SS, 0, 0xffffffff, 0xc093);
  xh_vm_set_register(0, VM_REG_GUEST_CS, 0x10);
  xh_vm_set_register(0, VM_REG_GUEST_DS, 0x18);
  xh_vm_set_register(0, VM_REG_GUEST_ES, 0x18);
  xh_vm_set_register(0, VM_REG_GUEST_SS, 0x18);
  
  // xh_vm_set_register(0, VM_REG_GUEST_RBP, 0);
  // xh_vm_set_register(0, VM_REG_GUEST_RDI, 0);
  // xh_vm_set_register(0, VM_REG_GUEST_RFLAGS, 0x2);
  // xh_vm_set_register(0, VM_REG_GUEST_RSI, 0x0);
  
  return boot_config->kernel_load_data.entry_addr;
}

uint64_t multiboot(void) {
  struct boot_config boot_config;

  FILE* kernel = fopen(config.kernel_path, "r");
  if(kernel == NULL) {
    perror("failed to open kernel");
    exit(1);
  }

  // peek at the first bit of the image, looking for a multiboot header,
  // if found, load any image specific config that we understand into 
  // boot_config
  if(0x0 == multiboot_find_header(kernel, &boot_config)) {
    printf("Didn't find a multiboot header in '%s'\n", config.kernel_path);
    exit(1);
  }

  // get the guest's memory range
  void* gpa = xh_vm_map_gpa(0, xh_vm_get_lowmem_size());
  boot_config.guest_mem_base = (uintptr_t)gpa;
  boot_config.guest_mem_size = xh_vm_get_lowmem_size();

  // actually load the image into the guest's memory
  size_t image_length = multiboot_load_image(kernel, &boot_config);
  if(0x0 == image_length) {
    printf("Failed to load kernel image into memory\n");
    exit(1);
  }

  // kernel image, we're done with you now.
  if(kernel != NULL)
    fclose(kernel);

  // load in all the specified modules
  char *s, *m;
  int mods_count = 0;

  s = m = config.module_list;
  
  if(config.module_list) {
    while(*m != 0x0) {
      while(*m != 0x0 && *m != ':') { m++; }
      if( *m == ':') m++;
      printf("module\n");
      mods_count++;
    }
  }






  return multiboot_set_guest_state(&boot_config);
}

//   // write out the multiboot info struct
//   void* p = (char*)((uintptr_t)host_load_addr + image_length);
//   struct multiboot_info* mb_info = (struct multiboot_info*)p;
//   p = (void*) ((uintptr_t)p +  sizeof(struct multiboot_info));
//   mb_info->flags = 0x0;


//   // write out all the modules!!
//   char *s, *m, *name, *v;

//   s = m = modulestring;


//   // count the number of modules
//   mb_info->mods_count = 0;
//   if(modulestring) {
//     while(*m != 0x0) {
//       while(*m != 0x0 && *m != ':') { m++; }
//       if( *m == ':') m++;
//       printf("module\n");
//       mb_info->mods_count++;
//     }
//   }
//   printf("There are %i modules\n", mb_info->mods_count);
  
//   struct multiboot_module_entry *table = (struct multiboot_module_entry*)p;
//   mb_info->mods_addr =(uint32_t)( (uintptr_t)p - (uintptr_t)gpa_map);

//   p = (void*) ((uintptr_t)p + sizeof(struct multiboot_module_entry) * mb_info->mods_count);
  
//   if(modulestring) {

//     mb_info->flags |= (1<<4);
//     printf("Writing out modules!\n");

//     s = m = modulestring;
//     while(*m != 0x0) {
//       while(*m != 0x0 && *m != ':') { m++; }
//       printf("p=%lu ",(uintptr_t) p);
//       p = ALIGN_4K(p);
//       printf("aligned p = %lu\n", (uintptr_t)p);
//       memcpy(p, s, (m - s) + 1);
//       name = p;
//       p =  (void*) ((uintptr_t)p + ((uintptr_t)m-(uintptr_t)s));
//       v = (char*)p;
//       *v = 0x0;
//       p =  (void*) ((uintptr_t)p + 1);
      
//       printf("Got a module: %s\n", name);

//       uint32_t module_size;

//       printf("Got module '%s'\n", name);

//       FILE* module = fopen(name, "r");
//       fseek(module, 0x0, SEEK_END);
//       module_size = (uint32_t)ftell(module);
//       fseek(module, 0x0, SEEK_SET);

//       printf("  size=%i bytes\n", module_size);

//       p = ALIGN_4K(p);
//       table->cmdline = (uint32_t)(name - (uintptr_t)gpa_map);
//       table->addr_start = (uint32_t)((uint32_t)p - (uintptr_t)gpa_map);
//       if( 1 != fread(p, module_size, 1, module)) perror("Failed to read module");
//       p = (void*)((uintptr_t)p + module_size);
//       table->addr_end = (uint32_t)((uint32_t)p - (uintptr_t)gpa_map);

//       fclose(module);

//       table++;

//       if( *m == ':') m++;
//       s = m;
//     }
//   }


//   if( boot_cmdline ) {
//     unsigned long length = (unsigned long)sprintf((char*)p,"%s %s", kernel_path_string, boot_cmdline);

//     mb_info->flags |= (1<<2);
//     mb_info->cmdline_addr = (uint32_t)((uintptr_t)p - (uintptr_t)gpa_map);
//     printf("cmdline addr is %x\n", mb_info->cmdline_addr);
//     p = (void*) ((uintptr_t)p +  length + 1);
//   }

//   mb_info->flags |= (1<<0);
//   mb_info->mem_lower = (uint32_t)640*1024;
//   mb_info->mem_upper = (uint32_t)xh_vm_get_lowmem_size();
