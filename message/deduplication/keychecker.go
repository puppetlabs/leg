package deduplication

// KeySetter is a type that can take a key and record its use. Keys here do not
// have associated values. This is not a key/value store, but rather a way to
// acknowledge the existence of a key. Keys are largely arbitrary strings, but
// should be passed into KeyValidators set by the implementation.
type KeySetter interface {
	// SetKeyAsSeen marks key as seen in the implementation's storage backend.
	// It is expected that any keys set remain set indefinitely and cannot be
	// removed (unless the entire storage backend is purged).
	SetKeyAsSeen(key string) error
}

// KeyChecker is a type that can take a key and report back information about
// it. This is not a key/value store, but rather a way to acknowledge the
// existence of a key.
type KeyChecker interface {
	// KeyHasBeenSeen checks if the key has been set in the implementation's
	// storage backend.
	KeyHasBeenSeen(key string) (bool, error)
}

// AtomicKeyCheckSetter is a type that can check if they key has been added
// before, then add it if it does not exist as an atomic operation.
type AtomicKeyCheckSetter interface {
	// CheckAndSetKey takes a key and checks if it exists. If does, then this
	// method returns true. If the key is not set and needs to be, then this
	// method returns false.
	CheckAndSetKey(key string) (bool, error)
}

// KeyValidator takes a key and validates it in some way. An example would be
// that the key is a certain legth or contains some data.
type KeyValidator interface {
	Apply(key string) error
}

// KeyValidatorFunc provides a generic implementation of KeyValidator.
type KeyValidatorFunc struct {
	Fn func(key string) error
}

func (k *KeyValidatorFunc) Apply(key string) error {
	return k.Fn(key)
}
