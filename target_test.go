package cross_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	type test struct {
		platform platform.Platform
	}
	tests := []test{
		{platform: platform.Windows},
		{platform: platform.Linux},
		{platform: platform.Darwin},
	}
	for _, test := range tests {
		t.Run(test.platform.String(), func(t *testing.T) {
			a := arch.AMD64
			target := cross.NewTest(test.platform, arch.AMD64)
			require.Equal(t, a, target.Architecture())
			require.Equal(t, test.platform, target.Platform())
		})
	}
}
