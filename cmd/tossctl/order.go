package main

import "github.com/spf13/cobra"

func newOrderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order",
		Short: "Preview and manage trading mutations",
		Long: "Trading commands are intentionally separate from read-only commands and " +
			"remain gated until trading discovery and permission controls are implemented.",
	}

	permissionsCmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage temporary trading execution permissions",
	}

	permissionsCmd.AddCommand(
		newStubCommand("grant", "Grant a short-lived trading permission"),
		newStubCommand("status", "Inspect the current trading permission state"),
		newStubCommand("revoke", "Revoke any active trading permission"),
	)

	cmd.AddCommand(
		newStubCommand("preview", "Preview a canonical order intent"),
		newStubCommand("place", "Place a live order with explicit danger approval"),
		newStubCommand("cancel", "Cancel a live pending order"),
		newStubCommand("amend", "Amend a live pending order"),
		permissionsCmd,
	)

	return cmd
}
