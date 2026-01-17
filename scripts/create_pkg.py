import os
import platform

def create_macos_pkg(script_content, pkg_name="PanelDemo", pkg_version="1.0"):
    """
    Creates a macOS package (.pkg) that executes a bash script.

    Args:
        script_content: The content of the bash script.
        pkg_name: The name of the package.
        pkg_version: The version of the package.

    Returns:
        The path to the created .pkg file if successful, None otherwise.
    """
    if platform.system() != "Darwin":
        print("This code is intended for macOS.")
        return None

    pkg_dir = os.path.expanduser("~")
    pkg_path = os.path.join(pkg_dir, f"{pkg_name}.pkg")

    with open(pkg_path, "w") as f:
        f.write(f"""# Package metadata
pkgName={pkg_name}
pkgVersion={pkg_version}

# Package contents
source = {script_content}
""")

    print(f"Package created at: {pkg_path}")
    print("To build the actual .pkg file, you need to use the `pkgbuild` command:")
    print(f"`sudo pkgbuild {pkg_path}`")
    print("Note: This is a simplified approach and might require further adjustments for a production-ready package.")

    return pkg_path

if __name__ == "__main__":
    bash_script = """#!/bin/bash

# Log file path
: "${LOG_FILE:=script.log}"

# Function to log a message with timestamp and arguments
log_message() {
  local timestamp=$(date +%Y-%m-%d %H:%M:%S)
  echo "[$timestamp] $* " >> "$LOG_FILE"
}
echo "$(date +%Y-%m-%d %H:%M:%S) : $*" >> /tmp/panel_log.log

log_message
echo "The package was installed."
"""

    pkg_path = create_macos_pkg(bash_script)
    if pkg_path:
        print(f"Successfully created package at: {pkg_path}")
        print(f"Now, build the package: sudo pkgbuild {pkg_path}")
    else:
        print("Failed to create package.")