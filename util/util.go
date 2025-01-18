package util

import (
	"regexp"
	pb "grpc-benchmark/protobuf"
)

const solanaAddressRegex = `^[1-9A-HJ-NP-Za-km-z]{32,44}$`

func IsValidSolanaAddress(address string) bool {
	re := regexp.MustCompile(solanaAddressRegex)
	return re.MatchString(address)
}

func BoolPtr(b bool) *bool {
	return &b
}
func CommitmentPtr(c pb.CommitmentLevel) * pb.CommitmentLevel {
	return &c
}