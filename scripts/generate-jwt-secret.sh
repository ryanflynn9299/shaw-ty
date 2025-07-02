#!/bin/bash

# Generate a secure random string for JWT secret
# Default length: 64 characters

LENGTH=${1:-64}

# Method 1: Using OpenSSL (most secure)
if command -v openssl &> /dev/null; then
    echo "Using OpenSSL to generate secret:"
    openssl rand -base64 $LENGTH | tr -d '\n'
    echo ""
    exit 0
fi

# Method 2: Using /dev/urandom
if [ -e /dev/urandom ]; then
    echo "Using /dev/urandom to generate secret:"
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c $LENGTH
    echo ""
    exit 0
fi

# Method 3: Fallback to less secure but widely available
echo "Using fallback method to generate secret:"
date +%s | sha256sum | base64 | head -c $LENGTH
echo ""