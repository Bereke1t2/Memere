import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_shadows.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_surface.dart';
import '../../../../shared/widgets/memere_mascot.dart';

class SavedScreen extends StatelessWidget {
  const SavedScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.bgPrimary,
      body: AppPageBackground(
        child: SafeArea(
          child: ListView(
            padding: const EdgeInsets.fromLTRB(
              AppSizes.screenPaddingH,
              AppSizes.lg,
              AppSizes.screenPaddingH,
              AppSizes.xl,
            ),
            children: [
              const Text('Saved', style: AppTextStyles.displayMedium),
              const SizedBox(height: AppSizes.xs),
              Text(
                'Bookmarked lessons and courses will appear here.',
                style: AppTextStyles.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: AppSizes.xl),
              AppSurface(
                padding: const EdgeInsets.all(AppSizes.lg),
                radius: AppSizes.radiusXl,
                color: AppColors.bgSecondary,
                shadows: AppShadows.md,
                child: Column(
                  children: [
                    const MemereMascot(
                      size: Size(220, 198),
                      showBackdrop: false,
                    ),
                    const SizedBox(height: AppSizes.lg),
                    const Text(
                      'No saved lessons yet',
                      textAlign: TextAlign.center,
                      style: AppTextStyles.headlineMedium,
                    ),
                    const SizedBox(height: AppSizes.sm),
                    Text(
                      'Use the bookmark button on courses you want to revisit.',
                      textAlign: TextAlign.center,
                      style: AppTextStyles.bodyMedium.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
