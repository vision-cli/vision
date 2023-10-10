$DOWNLOAD_URL = (curl https://api.github.com/repos/im2nguyen/rover/releases/latest | grep browser_download_url | grep darwin_arm64 | cut -d '"' -f 4)

$DOWNLOAD_ZIP = (curl --output-dir /tmp -OL https://github.com/im2nguyen/rover/releases/download/v0.3.3/rover_0.3.3_darwin_arm64.zip)

FUNC UNZIP = unzip /tmp/rover_0.3.3_darwin_arm64.zip -d /tmp/rover_0.3.3_darwin_arm64

FUNC MOVE = mv /tmp/rover_0.3.3_darwin_arm64/rover_0.3.3 ~/go/bin


