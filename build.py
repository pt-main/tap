"""
Universal Go cross-compilation script.

Builds a Go project for multiple operating systems and architectures.
All settings are configurable via command-line arguments.
"""

import argparse
import subprocess
import sys
import os
import shutil
from pathlib import Path
import re

# ----------------------------------------------------------------------
# Default values
DEFAULT_OUTPUT_DIR = "./build"
DEFAULT_NAME_TEMPLATE = "{project}-{os}-{arch}-{version}"
DEFAULT_PLATFORMS = [
    "windows/amd64",
    "windows/arm64",
    "linux/amd64",
    "linux/arm64",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
]
DEFAULT_VERSION = "final"

# ----------------------------------------------------------------------
def parse_platforms(platforms_str):
    """Parse a comma-separated list of 'os/arch' into a list of tuples."""
    if not platforms_str:
        return []
    items = [p.strip() for p in platforms_str.split(",") if p.strip()]
    result = []
    for item in items:
        parts = item.split("/")
        if len(parts) != 2:
            raise ValueError(f"Invalid platform format: '{item}'. Expected 'os/arch'.")
        result.append((parts[0], parts[1]))
    return result

def get_project_name(project_path):
    """Read the module name from go.mod to use as project name."""
    mod_file = Path(project_path) / "go.mod"
    if mod_file.exists():
        with open(mod_file, "r", encoding="utf-8") as f:
            for line in f:
                if line.startswith("module "):
                    # module name can be a path like github.com/user/project
                    # we take the last component as the project name
                    mod_name = line[len("module "):].strip()
                    # remove possible trailing comments
                    mod_name = mod_name.split()[0]
                    # return last part after slash, or the whole if no slash
                    return mod_name.split("/")[-1]
    # fallback to directory name
    return Path(project_path).name

def build_for_platform(project_path, output_dir, goos, goarch, name_template, version, ldflags, verbose):
    """Build the Go project for a single platform."""
    # Determine output file name
    # We replace placeholders in the template
    # Possible placeholders: {project}, {os}, {arch}, {version}
    # Also add .exe for windows if not already present in template
    project_name = get_project_name(project_path)
    name = name_template.format(
        project=project_name,
        os=goos,
        arch=goarch,
        version=version,
    )
    # Ensure .exe extension for windows if the template didn't already include it
    if goos == "windows" and not name.endswith(".exe"):
        name += ".exe"

    output_file = Path(output_dir) / name
    output_file.parent.mkdir(parents=True, exist_ok=True)

    # Build command
    cmd = ["go", "build"]
    if ldflags:
        cmd.extend(["-ldflags", ldflags])
    cmd.extend(["-o", str(output_file)])
    cmd.append(".")

    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch

    if verbose:
        print(f"\nBuilding for {goos}/{goarch} -> {output_file}")
        print(f"  Command: {' '.join(cmd)}")
        print(f"  GOOS={goos} GOARCH={goarch}")

    try:
        proc = subprocess.run(
            cmd,
            cwd=project_path,
            env=env,
            capture_output=True,
            text=True,
            check=False,
        )
        if proc.returncode != 0:
            print(f"ERROR: Build failed for {goos}/{goarch}")
            print(proc.stderr)
            return False
        if verbose:
            print(f"  Build successful.")
        return True
    except Exception as e:
        print(f"Exception while building for {goos}/{goarch}: {e}")
        return False

# ----------------------------------------------------------------------
def main():
    parser = argparse.ArgumentParser(
        description="Universal Go cross-compilation script",
        epilog="Example: %(prog)s -p ./myapp -o ./dist -n 'myapp-{os}-{arch}' -v 1.0.0 -pl windows/amd64,linux/amd64"
    )
    parser.add_argument(
        "-p", "--project-path",
        default=".",
        help="Path to the Go project (directory containing go.mod). Default: current directory."
    )
    parser.add_argument(
        "-o", "--output-dir",
        default=DEFAULT_OUTPUT_DIR,
        help=f"Output directory for binaries. Default: {DEFAULT_OUTPUT_DIR}"
    )
    parser.add_argument(
        "-n", "--name-template",
        default=DEFAULT_NAME_TEMPLATE,
        help=f"Name template for output files. Placeholders: {{project}}, {{os}}, {{arch}}, {{version}}. Default: '{DEFAULT_NAME_TEMPLATE}'"
    )
    parser.add_argument(
        "-pl", "--platforms",
        default=",".join(DEFAULT_PLATFORMS),
        help=f"Comma-separated list of platforms (os/arch). Default: {','.join(DEFAULT_PLATFORMS)}"
    )
    parser.add_argument(
        "-v", "--version",
        default=DEFAULT_VERSION,
        help=f"Version string to substitute in name template. Default: {DEFAULT_VERSION}"
    )
    parser.add_argument(
        "-ld", "--ldflags",
        default="",
        help="Additional ldflags to pass to go build (e.g. '-X main.version=1.0')"
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Print detailed build information."
    )

    args = parser.parse_args()

    # Resolve paths
    project_path = Path(args.project_path).resolve()
    if not project_path.exists():
        print(f"Error: Project path '{project_path}' does not exist.")
        sys.exit(1)

    output_dir = Path(args.output_dir).resolve()
    # Ensure output dir exists (will be created per platform)

    # Parse platforms
    try:
        platforms = parse_platforms(args.platforms)
    except ValueError as e:
        print(f"Error parsing platforms: {e}")
        sys.exit(1)

    if not platforms:
        print("No platforms specified. Nothing to build.")
        sys.exit(0)

    if args.verbose:
        print(f"Project path: {project_path}")
        print(f"Output directory: {output_dir}")
        print(f"Name template: {args.name_template}")
        print(f"Version: {args.version}")
        print(f"Platforms: {platforms}")
        print(f"Ldflags: {args.ldflags if args.ldflags else '(none)'}")

    # Build for each platform
    success_count = 0
    for goos, goarch in platforms:
        ok = build_for_platform(
            project_path=project_path,
            output_dir=output_dir,
            goos=goos,
            goarch=goarch,
            name_template=args.name_template,
            version=args.version,
            ldflags=args.ldflags,
            verbose=args.verbose,
        )
        if ok:
            success_count += 1

    print(f"\nBuild complete: {success_count}/{len(platforms)} succeeded.")
    if success_count < len(platforms):
        sys.exit(1)

if __name__ == "__main__":
    main()