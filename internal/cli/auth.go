package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Save your COROS credentials",
	Long: "Interactively choose your COROS region and enter your email and password. " +
		"The credentials are verified with a login, then saved: region and email in your user " +
		"config directory, and the password in the OS keyring.\n\n" +
		"If this fails, you can instead copy .env.example to .env and fill in your credentials.",
	Args: cobra.NoArgs,
}

const authFallbackHint = "Copy .env.example to .env and fill in your credentials to use them instead."

func init() {
	authCmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		out := cmd.OutOrStdout()
		in := cmd.InOrStdin()

		region, err := promptRegion(out, in)
		if err != nil {
			return authFailed(out, err)
		}

		email, err := promptEmail(out, in)
		if err != nil {
			return authFailed(out, err)
		}

		fmt.Fprint(out, "Password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(out)
		if err != nil {
			return authFailed(out, fmt.Errorf("reading password: %w", err))
		}
		password := string(passwordBytes)
		if password == "" {
			return authFailed(out, fmt.Errorf("password cannot be empty"))
		}

		// Verify the credentials before saving anything.
		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
		defer cancel()

		client := coros.NewClient(region.URL, email, password, "", http.DefaultClient)
		if err := client.Login(ctx); err != nil {
			return authFailed(out, fmt.Errorf("login failed: %w", err))
		}

		if err := config.SaveCredentials(region.URL, email, password); err != nil {
			return authFailed(out, err)
		}

		fmt.Fprintf(out, "Success: credentials for %s saved (%s region).\n", email, region.Name)
		return nil
	}
}

func authFailed(out io.Writer, err error) error {
	fmt.Fprintf(out, "Authentication failed: %v\n%s\n", err, authFallbackHint)
	return err
}

func promptRegion(out io.Writer, in io.Reader) (coros.Region, error) {
	names := make([]string, 0, len(coros.Regions))
	for _, r := range coros.Regions {
		names = append(names, r.Name)
	}

	const attempts = 3
	for i := 0; i < attempts; i++ {
		fmt.Fprintf(out, "Select your region (%s): ", strings.Join(names, ", "))

		var answer string
		if _, err := fmt.Fscanln(in, &answer); err != nil {
			return coros.Region{}, fmt.Errorf("reading region: %w", err)
		}

		region, err := coros.ParseRegion(answer)
		if err == nil {
			return region, nil
		}
		fmt.Fprintln(out, err)
	}

	return coros.Region{}, fmt.Errorf("no valid region selected")
}

func promptEmail(out io.Writer, in io.Reader) (string, error) {
	fmt.Fprint(out, "Email: ")

	var email string
	if _, err := fmt.Fscanln(in, &email); err != nil {
		return "", fmt.Errorf("reading email: %w", err)
	}

	email = strings.TrimSpace(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("invalid email %q", email)
	}

	return email, nil
}
