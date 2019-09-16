package maps

import "testing"

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "seeWhatIWant"}

	t.Run("Search for a known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "seeWhatIWant"
		assertString(t, got, want)
	})
	t.Run("Search for an unknown word", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		want := ErrNotFound.Error()
		assertError(t, err, ErrNotFound)
		assertString(t, err.Error(), want)
	})
}

func TestAdd(t *testing.T) {

	t.Run("Add new key to dict", func(t *testing.T) {
		dictionary := Dictionary{}
		err := dictionary.Add("test", "seeWhatIWant")

		got, err := dictionary.Search("test")
		want := "seeWhatIWant"
		assertNoError(t, err, "should find added word")
		assertString(t, got, want)
	})

	t.Run("Add word with existing key to dict ", func(t *testing.T) {
		dictionary := Dictionary{"test": "seeWhatIWant"}
		err := dictionary.Add("test", "newTest")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, "test", "seeWhatIWant")
	})
}

func TestUpdate(t *testing.T) {

	t.Run("Update existing key", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		newDefinition := "new definition"

		dictionary.Update(word, newDefinition)

		assertDefinition(t, dictionary, word, newDefinition)
	})

	t.Run("No key exist, try to update", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}

		err := dictionary.Update(word, definition)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete existing key", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}

		dictionary.Delete(word)
		_, err := dictionary.Search(word)
		if err != ErrNotFound {
			t.Errorf("Expected %q to be deleted", word)
		}
	})
}

func assertNoError(t *testing.T, got error, message string) {
	t.Helper()
	if got != nil {
		t.Fatal(message)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("expected to get an error.")
	}
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q for key: %q", got, want, "test")
	}
}

func assertDefinition(t *testing.T, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}
	if definition != got {
		t.Errorf("got %q want %q", got, definition)
	}
}
