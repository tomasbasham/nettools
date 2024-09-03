package cmd

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/tomasbasham/donut"
)

const maxBufferSize = 1024

func NewProxyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "proxy",
		Short: "Run a DNS proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			h := slog.NewJSONHandler(os.Stdout, nil)
			logger := slog.New(h)
			logger.Info("starting DNS proxy server")

			addr := &net.UDPAddr{
				IP:   net.IPv4(0, 0, 0, 0),
				Port: 53,
				Zone: "",
			}

			conn, err := net.ListenUDP("udp", addr)
			if err != nil {
				return err
			}
			defer conn.Close()

			// Create a context that will be canceled when a signal is received.
			// This allows us to gracefully shutdown the server when a signal is
			// received.
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			buf := make([]byte, maxBufferSize)

			// Given that waiting for packets to arrive is blocking by nature and we
			// want to be able of canceling such action if desired, we do that in a
			// separate go routine.
			go func() {
				for {
					n, addr, err := conn.ReadFromUDP(buf)
					if err != nil {
						logger.Error("failed to read from UDP connection: " + err.Error())
						continue
					}

					// handle the request
					go handleRequest(conn, addr, buf[:n])
				}
			}()

			<-ctx.Done()
			return nil
		},
	}
}

func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	resolver := donut.New(donut.GoogleHost)
	message, err := resolver.LookupRaw(buf)
	if err != nil {
		panic(err)
	}

	b, err := conn.WriteTo(message, addr)
	if err != nil {
		panic(err)
	}

	if b != len(message) {
		panic("message not sent")
	}
}
