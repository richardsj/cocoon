# Run an API. 

repoParent="/go/src/github.com/ncodes"
mkdir -p $repoParent
cd $repoParent

# pull cocoon source
git clone --depth=1 -b connector-redesign https://github.com/ncodes/cocoon

# start the orderer
cd cocoon
glide install
go run core/main.go api start
