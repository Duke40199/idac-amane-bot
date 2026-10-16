// Package port defines the interfaces (ports) through which the domain talks to
// adapters: repositories, notifiers, external clients.
//
// It is intentionally empty in this wireframe: all IDAC-specific ports have been
// removed. Add interfaces here as features grow; adapters implement them and the
// domain depends only on these interfaces.
package port
