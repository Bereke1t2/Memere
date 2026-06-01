#!/bin/bash
# Core folders
mkdir -p lib/core/constants
mkdir -p lib/core/errors
mkdir -p lib/core/network
mkdir -p lib/core/storage
mkdir -p lib/core/theme
mkdir -p lib/core/router
mkdir -p lib/core/di

# Auth feature
mkdir -p lib/features/auth/data/datasources
mkdir -p lib/features/auth/data/models
mkdir -p lib/features/auth/data/repositories
mkdir -p lib/features/auth/domain/entities
mkdir -p lib/features/auth/domain/repositories
mkdir -p lib/features/auth/domain/usecases
mkdir -p lib/features/auth/presentation/providers
mkdir -p lib/features/auth/presentation/screens
mkdir -p lib/features/auth/presentation/widgets

# Stub folders for future phases
mkdir -p lib/features/courses/data/datasources
mkdir -p lib/features/courses/data/models
mkdir -p lib/features/courses/data/repositories
mkdir -p lib/features/courses/domain/entities
mkdir -p lib/features/courses/domain/repositories
mkdir -p lib/features/courses/domain/usecases
mkdir -p lib/features/courses/presentation/providers
mkdir -p lib/features/courses/presentation/screens
mkdir -p lib/features/courses/presentation/widgets

mkdir -p lib/features/video_player/presentation/screens
mkdir -p lib/features/video_player/presentation/widgets
mkdir -p lib/features/quiz/presentation/screens
mkdir -p lib/features/exam/presentation/screens
mkdir -p lib/features/payment/presentation/screens
mkdir -p lib/features/progress/presentation/screens
mkdir -p lib/features/notifications/presentation/screens

# Shared
mkdir -p lib/shared/widgets
mkdir -p lib/shared/extensions
mkdir -p lib/shared/utils

# Assets
mkdir -p assets/images
mkdir -p assets/icons
mkdir -p assets/lottie
mkdir -p assets/fonts

# Tests
mkdir -p test/unit/features/auth
mkdir -p test/widget/features/auth
mkdir -p test/integration

# Docs (already exists)
mkdir -p docs
