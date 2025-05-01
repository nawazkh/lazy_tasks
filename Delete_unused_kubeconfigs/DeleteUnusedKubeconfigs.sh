#!/bin/bash

# Create the backup directory if it doesn't exist
backup_dir=~/.kube/backups
mkdir -p "$backup_dir"

# Get the current timestamp in a human-readable format with AM/PM
timestamp=$(date +"%Y-%m-%d_%I-%M-%S_%p")

# Set the kubeconfig file path
KUBECONFIG_FILE=~/.kube/config

# Backup the original kubeconfig file with the timestamp into the backup directory
backup_file="$backup_dir/kubeconfig_backup_$timestamp"
cp $KUBECONFIG_FILE "$backup_file"
echo -e "\n\033[1;34mBackup of kubeconfig created at $backup_file\033[0m\n"

# Get the current context
current_context=$(kubectl config current-context)
echo -e "\n\033[1;32mCurrent context: $current_context\033[0m\n"

# List all kubeconfig contexts except the current one
unused_contexts=$(kubectl config get-contexts -o name | grep -v "$current_context")

# Delete all unused contexts
echo -e "\033[1;33m--- Starting Context Cleanup ---\033[0m"
for context in $unused_contexts; do
    echo -e "\033[1;31mRemoving unused context: $context\033[0m"
    kubectl config delete-context "$context" > /dev/null
done
echo -e "\033[1;33m--- Context Cleanup Completed ---\033[0m\n"

# Get all clusters in use by current contexts
clusters_in_use=$(kubectl config view -o jsonpath='{.contexts[*].context.cluster}')

# Get all users in use by current contexts
users_in_use=$(kubectl config view -o jsonpath='{.contexts[*].context.user}')

# Unset unused clusters
echo -e "\033[1;33m--- Starting Cluster Cleanup ---\033[0m"
for cluster in $(kubectl config view -o jsonpath='{.clusters[*].name}'); do
    if ! echo "$clusters_in_use" | grep -q "$cluster"; then
        echo -e "\033[1;31mRemoving unused cluster: $cluster\033[0m"
        kubectl config unset clusters."$cluster" > /dev/null
    fi
done
echo -e "\033[1;33m--- Cluster Cleanup Completed ---\033[0m\n"


# Unset unused users
echo -e "\033[1;33m--- Starting User Cleanup ---\033[0m"
for user in $(kubectl config view -o jsonpath='{.users[*].name}'); do
    if ! echo "$users_in_use" | grep -q "$user"; then
        echo -e "\033[1;31mRemoving unused user: $user\033[0m"
        kubectl config unset users."$user" > /dev/null
    fi
done
echo -e "\033[1;33m--- User Cleanup Completed ---\033[0m\n"

# Final message
echo -e "\033[1;32m--- All cleanups completed successfully! Current context: $current_context ---\033[0m"
