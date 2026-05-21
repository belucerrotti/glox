import sys
import subprocess
import os

# busca el binario: primero el local ./glox, sino en el PATH como "glox"
script_dir = os.path.dirname(os.path.abspath(__file__))
local_bin = os.path.join(script_dir, "..", "glox")
LOX_BINARY = local_bin if os.path.isfile(local_bin) else "glox"

for lox_file in filter(lambda f: f.endswith(".lox"), sorted(os.listdir("real-tests"))):
    print(f"$ {LOX_BINARY} real-tests/{lox_file}")
    result = subprocess.run(
        [LOX_BINARY, f"real-tests/{lox_file}"],
        stdout=subprocess.PIPE,
        stdin=subprocess.DEVNULL,
    )
    out = result.stdout.decode().strip()
    print(out)
    print()
    if "ERROR".lower() in out.lower():
        sys.exit(1)

print("Todo OK")
