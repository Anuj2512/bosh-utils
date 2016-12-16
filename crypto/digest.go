package crypto

import (
	"errors"
	"fmt"
	"io"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DigestAlgorithm string

type digestImpl struct {
	algorithm Algorithm
	digest    string
}

type algorithmSHA512 struct {
	Name string
}

func (a algorithmSHA512)  String() string {
	return a.Name
}

func (a algorithmSHA512)  Compare(other Algorithm) int {
	if _, ok := other.(algorithmSHA512); ok {
		return 0
	}
	return 1
}

func (a algorithmSHA512)  CreateDigest(reader io.Reader) (Digest, error) {
	return digestProviderImpl{}.CreateFromStream(reader, DigestAlgorithm(a.Name))
}

type algorithmSHA256 struct {
	Name string
}

func (a algorithmSHA256)  CreateDigest(reader io.Reader) (Digest, error) {
	return digestProviderImpl{}.CreateFromStream(reader, DigestAlgorithm(a.Name))

}

func (a algorithmSHA256)  String() string {
	return a.Name
}

func (a algorithmSHA256)  Compare(other Algorithm) int {
	if _, ok := other.(algorithmSHA1); ok {
		return 1
	} else if _, ok = other.(algorithmSHA512); ok {
		return -1
	} else if _, ok = other.(algorithmSHA256); ok {
		return 0
	}
	return 1
}

type algorithmSHA1 struct {
	Name string
}

func (a algorithmSHA1)  CreateDigest(reader io.Reader) (Digest, error) {
	return digestProviderImpl{}.CreateFromStream(reader, DigestAlgorithm(a.Name))

}

func (a algorithmSHA1)  String() string {
	return a.Name
}

func (a algorithmSHA1)  Compare(other Algorithm) int {
	if _, ok := other.(algorithmSHA1); ok {
		return 0
	} else if _, ok = other.(algorithmSHA512); ok {
		return -1
	} else if _, ok = other.(algorithmSHA256); ok {
		return -1
	}
	return 1
}

func (c digestImpl) Algorithm() Algorithm {
	return c.algorithm
}

func (c digestImpl) String() string {
	if c.algorithm == DigestAlgorithmSHA1 {
		return c.digest
	}

	return fmt.Sprintf("%s:%s", c.algorithm, c.digest)
}

func (c digestImpl) Verify(reader io.Reader) error {
	otherDigest, err := c.Algorithm().CreateDigest(reader)
	if err != nil {
		return err
	}

	if otherDigest, ok := otherDigest.(digestImpl); ok {
		if c.digest != otherDigest.digest {
			return errors.New(fmt.Sprintf(`Expected %s digest "%s" but received "%s"`, c.algorithm, c.digest, otherDigest.String()))
		}
	} else {
		return errors.New(fmt.Sprintf(`Unknown digest to verify against %v`, otherDigest.String()))
	}

	return nil
}

func NewAlgorithm(algorithm string) (Algorithm, error) {
	switch algorithm {
	case "sha1":
		return DigestAlgorithmSHA1, nil
	case "sha256":
		return DigestAlgorithmSHA256, nil
	case "sha512":
		return DigestAlgorithmSHA512, nil
	}
	return nil, bosherr.Errorf("Unknown algorithim '%s'", algorithm)
}

func NewDigest(algorithm Algorithm, digest string) digestImpl {
	return digestImpl{
		algorithm: algorithm,
		digest:    digest,
	}
}
