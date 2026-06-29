import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';

class QuizResultSkeleton extends StatelessWidget {
  const QuizResultSkeleton({super.key});

  @override
  Widget build(BuildContext context) {
    return Shimmer.fromColors(
      baseColor: AppColors.bgSecondary,
      highlightColor: AppColors.bgTertiary,
      child: ListView(
        padding: const EdgeInsets.all(AppSizes.screenPaddingH),
        children: const [
          _Block(height: 190),
          SizedBox(height: AppSizes.md),
          _Block(height: 120),
          SizedBox(height: AppSizes.md),
          _Block(height: 88),
          SizedBox(height: AppSizes.sm),
          _Block(height: 88),
          SizedBox(height: AppSizes.sm),
          _Block(height: 88),
          SizedBox(height: AppSizes.sm),
          _Block(height: 88),
        ],
      ),
    );
  }
}

class _Block extends StatelessWidget {
  const _Block({required this.height});

  final double height;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      ),
    );
  }
}
