import sys
import subprocess
import os

GREEN = "\033[92m"
RED   = "\033[91m"
CYAN  = "\033[96m"
BOLD  = "\033[1m"
RESET = "\033[0m"
DIM   = "\033[2m"

script_dir = os.path.dirname(os.path.abspath(__file__))
local_bin = os.path.join(script_dir, "..", "glox")
LOX_BINARY = local_bin if os.path.isfile(local_bin) else "glox"

total_ok = 0
total_fail = 0

def run_suite(folder):
    global total_ok, total_fail
    path = os.path.join(script_dir, folder)
    files = sorted(f for f in os.listdir(path) if f.endswith(".lox"))
    print(f"\n{BOLD}{CYAN}{'─'*40}{RESET}")
    print(f"{BOLD}{CYAN}  {folder.upper()}  ({len(files)} archivos){RESET}")
    print(f"{BOLD}{CYAN}{'─'*40}{RESET}")

    for lox_file in files:
        full_path = os.path.join(path, lox_file)
        result = subprocess.run(
            [LOX_BINARY, full_path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=subprocess.DEVNULL,
        )
        out = result.stdout.decode().strip()
        err = result.stderr.decode().strip()
        failed = "error" in out.lower() or result.returncode != 0

        name = lox_file.replace(".lox", "")
        status = f"{RED}✗ FAIL{RESET}" if failed else f"{GREEN}✓ OK{RESET}"
        print(f"  {status}  {name}")

        # imprimir las líneas del output indentadas
        for line in out.splitlines():
            print(f"    > {line}")
        if err:
            print(f"    [stderr] {err}")

        if failed:
            total_fail += 1
        else:
            total_ok += 1

run_suite("basic")
run_suite("advanced")

print(f"\n{BOLD}{'─'*40}{RESET}")
total = total_ok + total_fail
if total_fail == 0:
    print(f"{BOLD}{GREEN}  TODOS LOS TESTS PASARON ({total_ok}/{total}){RESET}")
else:
    print(f"{BOLD}{RED}  {total_fail} TEST(S) FALLARON  ({total_ok}/{total} OK){RESET}")
print(f"{BOLD}{'─'*40}{RESET}\n")

if total_fail > 0:
    sys.exit(1)
