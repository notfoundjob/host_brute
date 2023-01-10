# X86
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/host_brute_linux_amd64 host_brute.go
# ARM
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/host_brute_linux_arm64 host_brute.go


# X86
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/host_brute_windows_amd64 host_brute.go
 
# ARM
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./release/host_brute_windows_arm64 host_brute.go

# X86
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/host_brute_darwin_amd64 host_brute.go
 
# ARM
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./release/host_brute_darwin_arm64 host_brute.go