#!/bin/sh
# Fetch and patch the SkyWater Sky130 primitive device models for the
# Tier B SPICE runs. The library lands in tierb/sky130_fd_pr (gitignored,
# ~770 MB). Two fixes make the raw library ngspice-compatible:
#   1. spectre-style 'include' lines become '.include'
#   2. one ESD model is referenced by a filename that does not ship
set -e
cd "$(dirname "$0")"
if [ ! -d sky130_fd_pr ]; then
    git clone --depth 1 \
        https://github.com/google/skywater-pdk-libs-sky130_fd_pr.git \
        sky130_fd_pr
fi
grep -rl '^include ' sky130_fd_pr/cells sky130_fd_pr/models 2>/dev/null \
    | xargs -r sed -i 's/^include /.include /'
cp -f sky130_fd_pr/cells/esd_nfet_05v0_nvt/sky130_fd_pr__esd_nfet_05v0_nvt.pm3.spice \
      sky130_fd_pr/cells/nfet_05v0_nvt/sky130_fd_pr__esd_nfet_05v0_nvt.pm3
echo "Sky130 models ready."
