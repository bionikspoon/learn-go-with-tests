package maps

import "testing"

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		assertDefinitions(t, dictionary, "test", "this is just a test")
	})

	t.Run("unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")

		assertError(t, err, ErrNotFound)
	})

}

func TestAdd(t *testing.T) {

	t.Run("new word", func(t *testing.T) {
		dictionary := make(Dictionary)

		err := dictionary.Add("test", "this is just a test")

		assertError(t, err, nil)
		assertDefinitions(t, dictionary, "test", "this is just a test")
	})

	t.Run("existing word", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is just a test"}

		err := dictionary.Add("test", "this value is not used")
		assertError(t, err, ErrWordExists)

		assertDefinitions(t, dictionary, "test", "this is just a test")
	})
}

func TestUpdate(t *testing.T) {

	t.Run("new word", func(t *testing.T) {
		dictionary := make(Dictionary)

		err := dictionary.Update("test", "this is just a test")

		assertError(t, err, ErrWordDoesNotExist)

	})

	t.Run("existing word", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is just a test"}

		err := dictionary.Update("test", "new definition")

		assertError(t, err, nil)
		assertDefinitions(t, dictionary, "test", "new definition")
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		dictionary := Dictionary{"test": "this is just a test"}

		err := dictionary.Delete("test")

		assertError(t, err, nil)

		_, err = dictionary.Search("test")

		assertError(t, err, ErrNotFound)
	})

	t.Run("missing word", func(t *testing.T) {
		dictionary := make(Dictionary)

		err := dictionary.Delete("test")

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func assertDefinitions(t *testing.T, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)

	assertError(t, err, nil)
	assertStrings(t, got, definition)
}

func assertStrings(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
