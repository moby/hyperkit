(* Lightweight thread library for Objective Caml
 * http://www.ocsigen.org/lwt
 * Program discover
 * Copyright (C) 2010 Jérémie Dimino
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as
 * published by the Free Software Foundation, with linking exceptions;
 * either version 2.1 of the License, or (at your option) any later
 * version. See COPYING file for details.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA
 * 02111-1307, USA.
 *)

(* Discover available features *)

let cut_tail l = List.rev (List.tl (List.rev l))

let string_split sep source =
  let copy_part index offset =
    String.sub source index (offset - index)
  in
  let l = String.length source in
  let rec loop prev current acc =
    if current >= l then
      List.rev acc
    else
      match (source.[current] = sep, current = prev, current = l - 1) with
      | (true, true, _) -> loop (current + 1) (current + 1) acc
      | (true, _, _) -> loop (current + 1) (current + 1) ((copy_part prev current)::acc)
      | (false, _, true) -> loop (current + 1) (current + 1) ((copy_part prev (current + 1))::acc)
      | _ -> loop prev (current + 1) acc
  in loop 0 0 []

let uniq lst =
  let unique_set = Hashtbl.create (List.length lst) in
  List.iter (fun x -> Hashtbl.replace unique_set x ()) lst;
  Hashtbl.fold (fun x () xs -> x :: xs) unique_set []

let get_paths env_name =
  try
    let paths = Sys.getenv env_name in
    let dirs = string_split ':' paths in
    List.map (fun dir ->
      let components = string_split '/' dir in
      "/" ^ (String.concat "/" (cut_tail components))
    ) dirs
  with Not_found -> []

let env_paths = List.append (get_paths "LIBRARY_PATH") (get_paths "C_INCLUDE_PATH")

(* Keep that in sync with the list in myocamlbuild.ml *)
let search_paths = uniq (List.append [
  "/usr";
  "/usr/local";
  "/opt";
  "/opt/local";
  "/sw";
  "/mingw";] env_paths)

open Printf

(* +-----------------------------------------------------------------+
   | Test codes                                                      |
   +-----------------------------------------------------------------+ *)

let caml_code = "
external test : unit -> unit = \"lwt_test\"
let () = test ()
"


let libev_code = "
#include <caml/mlvalues.h>
#include <ev.h>

CAMLprim value lwt_test()
{
  ev_default_loop(0);
  return Val_unit;
}
"
(* +-----------------------------------------------------------------+
   | Compilation                                                     |
   +-----------------------------------------------------------------+ *)

let ocamlc = ref "ocamlc"
let ext_obj = ref ".o"
let exec_name = ref "a.out"

let log_file = ref ""
let caml_file = ref ""

(* Search for a header file in standard directories. *)
let search_header header =
  let rec loop = function
    | [] ->
        None
    | dir :: dirs ->
        if Sys.file_exists (dir ^ "/include/" ^ header) then
          Some dir
        else
          loop dirs
  in
  loop search_paths

let c_args =
  let flags path = Printf.sprintf "-ccopt -I%s/include -ccopt -L%s/lib" path path in
  match search_header "ev.h" with
    | None -> ""
    | Some path -> flags path

let compile args stub_file =
  let cmd = sprintf "%s -custom %s %s %s %s > %s 2>&1"
    !ocamlc
    c_args
    (Filename.quote stub_file)
    args
    (Filename.quote !caml_file)
    (Filename.quote !log_file) in
  Sys.command cmd = 0

let safe_remove file_name =
  try
    Sys.remove file_name
  with exn ->
    ()

let test_code args stub_code =
  let stub_file, oc = Filename.open_temp_file "lwt_stub" ".c" in
  let cleanup () =
    safe_remove stub_file;
    safe_remove (Filename.chop_extension (Filename.basename stub_file) ^ !ext_obj)
  in
  try
    output_string oc stub_code;
    flush oc;
    close_out oc;
    let result = compile args stub_file in
    cleanup ();
    result
  with exn ->
    (try close_out oc with _ -> ());
    cleanup ();
    raise exn

let config = open_out "lwt_config.h"
let config_ml = open_out "lwt_config.ml"

let test_feature ?(do_check = true) name macro ?(args="") code =
  if do_check then begin
    printf "testing for %s:%!" name;
    if test_code args code then begin
      fprintf config "#define %s\n" macro;
      fprintf config_ml "#let %s = true\n" macro;
      printf " %s available\n%!" (String.make (34 - String.length name) '.');
      true
    end else begin
      fprintf config "//#define %s\n" macro;
      fprintf config_ml "#let %s = false\n" macro;
      printf " %s unavailable\n%!" (String.make (34 - String.length name) '.');
      false
    end
  end else begin
    printf "not checking for %s\n%!" name;
    fprintf config "//#define %s\n" macro;
    fprintf config_ml "#let %s = false\n" macro;
    true
  end

(* +-----------------------------------------------------------------+
   | Entry point                                                     |
   +-----------------------------------------------------------------+ *)

let () =
  let args = [
    "-ocamlc", Arg.Set_string ocamlc, "<path> ocamlc";
    "-ext-obj", Arg.Set_string ext_obj, "<ext> C object files extension";
    "-exec-name", Arg.Set_string exec_name, "<name> name of the executable produced by ocamlc";
  ] in
  Arg.parse args ignore "check for external C libraries and available features\noptions are:";

  (* Put the caml code into a temporary file. *)
  let file, oc = Filename.open_temp_file "lwt_caml" ".ml" in
  caml_file := file;
  output_string oc caml_code;
  close_out oc;

  log_file := Filename.temp_file "lwt_output" ".log";

  (* Cleanup things on exit. *)
  at_exit (fun () ->
             (try close_out config with _ -> ());
             (try close_out config_ml with _ -> ());
             safe_remove !log_file;
             safe_remove !exec_name;
             safe_remove !caml_file;
             safe_remove (Filename.chop_extension !caml_file ^ ".cmi");
             safe_remove (Filename.chop_extension !caml_file ^ ".cmo"));

  let missing = [] in
  let missing = if test_feature "libev" "HAVE_LIBEV" ~args:"-cclib -lev" libev_code then missing else "libev" :: missing in

  if missing <> [] then begin
    printf "
      The following recquired C libraries are missing: %s.
Please install them and retry. If they are installed in a non-standard location, set the environment variables C_INCLUDE_PATH and LIBRARY_PATH accordingly and retry.

For example, if they are installed in /opt/local, you can type:

export C_INCLUDE_PATH=/opt/local/include
export LIBRARY_PATH=/opt/local/lib

To compile without libev support, use ./configure --disable-libev ...
" (String.concat ", " missing);
    exit 1
  end;

(*
  ignore (test_feature "eventfd" "HAVE_EVENTFD" eventfd_code);
  ignore (test_feature "fd passing" "HAVE_FD_PASSING" fd_passing_code);
  ignore (test_feature "sched_getcpu" "HAVE_GETCPU" getcpu_code);
  ignore (test_feature "affinity getting/setting" "HAVE_AFFINITY" affinity_code);
  ignore (test_feature "credentials getting" "HAVE_GET_CREDENTIALS" get_credentials_code);
  ignore (test_feature "fdatasync" "HAVE_FDATASYNC" fdatasync_code)
*)
