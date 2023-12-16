rm -rf /tmp/underbeing

rm -f /tmp/underbeing.tar
rm -f /tmp/filelist.txt

{
    rg --files . \
        | grep -vE underbeing$ \
        | grep -v README.org \
        | grep -v make_txtar.sh \
        | grep -v go.sum \
        | grep -v go.mod \
        | grep -v Makefile \
        | grep -v options/options.go \
        | grep -v cmd/main.go \
        # | grep -v underbeing.go \

} | tee /tmp/filelist.txt
tar -cf /tmp/underbeing.tar -T /tmp/filelist.txt
mkdir -p /tmp/underbeing
tar xf /tmp/underbeing.tar -C /tmp/underbeing
rg --files /tmp/underbeing
txtar-c /tmp/underbeing | pbcopy
