# Delete Unused Kubeconfigs

A bash script that safely cleans up your Kubernetes configuration by removing unused contexts, clusters, and users while preserving your current active context.

## Overview

This script helps maintain a clean Kubernetes configuration by:
1. Creating a backup of your current kubeconfig
2. Identifying and removing unused contexts
3. Cleaning up unused cluster configurations
4. Removing unused user credentials
5. Preserving your current active context throughout the process

## Features

- **Automatic Backup**: Creates timestamped backups of your kubeconfig in `~/.kube/backups`
- **Safe Execution**: Preserves your current context and its dependencies
- **Comprehensive Cleanup**: Removes:
  - Unused contexts
  - Unused cluster configurations
  - Unused user credentials
- **Color-coded Output**: Provides clear visual feedback during execution

## Prerequisites

- Kubernetes CLI (`kubectl`) installed and configured
- Active kubeconfig file at `~/.kube/config`

## Usage

```bash
./DeleteUnusedKubeconfigs.sh
```

## What it Does

1. **Backup Creation**
   - Creates a backup directory at `~/.kube/backups`
   - Saves a timestamped copy of your current kubeconfig

2. **Context Management**
   - Identifies your current active context
   - Removes all other unused contexts

3. **Cluster Cleanup**
   - Identifies clusters referenced by active contexts
   - Removes cluster configurations that aren't in use

4. **User Cleanup**
   - Identifies users referenced by active contexts
   - Removes user credentials that aren't in use

## Output

The script provides clear, color-coded output indicating:
- Location of the backup file
- Current active context
- Lists of removed contexts, clusters, and users
- Confirmation of successful completion

## Safety Features

- Automatic backup creation before any changes
- Preservation of current context and its dependencies
- Non-destructive to active configurations

## Note

Always ensure you have necessary permissions and valid backup before running cleanup scripts on your Kubernetes configuration.

## Sample output
```bash
❯ delKConfigs=/Users/nawazhussain/msftcode/DeleteUnusedKubeconfigs.sh
❯ delKConfigs

Backup of kubeconfig created at /Users/nawazhussain/.kube/backups/kubeconfig_backup_2025-05-01_03-22-49_PM


Current context: aks-mgmt-1745538551

--- Starting Context Cleanup ---
Removing unused context: aks-mgmt-1744307866
.
.
.
Removing unused context: capz-e2e-9
--- Context Cleanup Completed ---

--- Starting Cluster Cleanup ---
Removing unused cluster: aks-mgmt-1744307866
.
.
.
--- Cluster Cleanup Completed ---

--- Starting User Cleanup ---
Removing unused user: clusterUser_aks-mgmt-1744307866_aks-mgmt-1744307866
.
.
.
--- User Cleanup Completed ---

--- All cleanups completed successfully! Current context: aks-mgmt-1745538551 ---
```