package maps

type Dictionary map[string]string

const (
	ErrNotFound         = DictionaryError("could not find the word you were looking for")
	ErrWordExists       = DictionaryError("could not add word because it already exists")
	ErrWordDoesNotExist = DictionaryError("could not update/delete word because it does not exist")
)

type DictionaryError string

func (err DictionaryError) Error() string {
	return string(err)
}

func (dictionary Dictionary) Search(word string) (string, error) {
	definition, ok := dictionary[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}

func (dictionary Dictionary) Add(word, definition string) error {
	_, err := dictionary.Search(word)

	switch err {
	case ErrNotFound:
		dictionary[word] = definition
		return nil
	case nil:
		return ErrWordExists
	default:
		return err
	}
}

func (dictionary Dictionary) Update(word, definition string) error {
	_, err := dictionary.Search(word)

	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		dictionary[word] = definition
		return nil
	default:
		return err
	}
}

func (dictionary Dictionary) Delete(word string) error {
	_, err := dictionary.Search(word)

	switch err {
	case nil:
		delete(dictionary, word)
		return nil
	case ErrNotFound:
		return ErrWordDoesNotExist
	default:
		return err
	}

}
