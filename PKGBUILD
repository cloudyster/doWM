# Maintainer: BobdaProgrammer <https://github.com/BobdaProgrammer>
pkgname=doWM
pkgver=1.0.3
pkgrel=1
pkgdesc="A beautiful tiling and floating x11 window manager"
arch=('x86_64')
url="https://github.com/BobdaProgrammer/doWM"
license=('MIT')
depends=('xorg-server')
makedepends=('go' 'git')
source=("doWM.desktop"
    "LICENSE")
sha256sums=('SKIP' 'SKIP')

build() {
    cd "$srcdir"
    cd ..
    go build -o doWM
}

package() {
  # Install binary
  sudo install -Dm755 "../doWM" "/usr/local/bin/doWM"

  # Install .desktop session file
  sudo install -Dm644 "doWM.desktop" "/usr/share/xsessions/doWM.desktop"

  # License
  sudo install -vDm644 "LICENSE" "/usr/share/licenses/$pkgname" 
}
