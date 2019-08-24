package pointer

import (
	"testing"
)

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(tt *testing.T) {

		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))

		assertBalance(tt, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(tt *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(tt, wallet, Bitcoin(10))
		assertNoError(tt, err)
	})

	t.Run("Withdraw insufficient funds", func(tt *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(tt, wallet, startingBalance)
		assertError(tt, err, ErrInsufficientFunds)

	})
}

func TestBitcoin_String(t *testing.T) {
	tests := []struct {
		name string
		b    Bitcoin
		want string
	}{
		{"Bitcoin Stringer", Bitcoin(11), "11 BTC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Bitcoin.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("expecting an error")
	}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	t.Helper()

	if got := wallet.Balance(); got != want {
		t.Errorf("%#v got %s want %s", wallet, got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatal("not expecting an error")
	}
}
