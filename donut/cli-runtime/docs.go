/**
 * A set of helpers to make it easier to build CLI applications in Go. The
 * package is a shared dependency for clients to work with StatusCake
 * infrastructure and services that allows to maintain compatible behaviour
 * across applications.
 *
 * This package should not be extended to include any application specific types
 * of helper funcitons, neither should it include API types. These should
 * instead be implemented in the application specific package.
 */

package cliruntime // import "github.com/tomasbasham/donut/cli-runtime"
