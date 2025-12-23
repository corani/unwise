package dummy

// NOTE(daniel): this causes `mockery` to berun for the top-level project. Using the
// `.mockery.yaml` config file, it'll generate mocks for all required interfaces in
// the project.
//
//go:generate go tool mockery
