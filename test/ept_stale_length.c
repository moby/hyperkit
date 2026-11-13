/* SPDX-License-Identifier: BSD-2-Clause */

/* Supply defined test values for an architecturally undefined MMIO length. */
#include <Hypervisor/hv.h>
#include <Hypervisor/hv_vmx.h>
#include <Hypervisor/hv_arch_vmx.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

static unsigned int instruction_length;

static void __attribute__((constructor))
configure(void)
{
	const char *setting;
	char *end;
	unsigned long value;

	setting = getenv("HYPERKIT_TEST_EPT_LENGTH");
	if (setting == NULL)
		abort();
	value = strtoul(setting, &end, 10);
	if (*setting == '\0' || *end != '\0' || value > 15)
		abort();
	instruction_length = (unsigned int)value;
}

static hv_return_t
stale_length(hv_vcpuid_t cpu, uint32_t field, uint64_t *value)
{
	hv_return_t result;
	uint64_t reason, gpa;

	result = hv_vmx_vcpu_read_vmcs(cpu, field, value);
	if (result == HV_SUCCESS && field == VMCS_RO_VMEXIT_INSTR_LEN &&
	    hv_vmx_vcpu_read_vmcs(cpu, VMCS_RO_EXIT_REASON, &reason) == HV_SUCCESS &&
	    reason == 48 &&
	    hv_vmx_vcpu_read_vmcs(cpu, VMCS_GUEST_PHYSICAL_ADDRESS, &gpa) == HV_SUCCESS &&
	    gpa == 0xfec00000) {
		*value = instruction_length;
		fprintf(stderr, "EPT_TEST_LENGTH=%u\n", instruction_length);
	}
	return (result);
}

__attribute__((used, section("__DATA,__interpose")))
static const struct {
	const void *replacement;
	const void *original;
} interpose = {
	(const void *)&stale_length,
	(const void *)&hv_vmx_vcpu_read_vmcs
};
