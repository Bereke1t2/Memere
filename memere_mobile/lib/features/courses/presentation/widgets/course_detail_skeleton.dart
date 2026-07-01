import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';

class CourseDetailSkeleton extends StatelessWidget {
  const CourseDetailSkeleton({super.key});

  @override
  Widget build(BuildContext context) {
    return Shimmer.fromColors(
      baseColor: AppColors.skeletonBase,
      highlightColor: AppColors.skeletonHighlight,
      child: ListView(
        padding: const EdgeInsets.all(AppSizes.screenPaddingH),
        children: [
          _block(height: 188),
          const SizedBox(height: AppSizes.md),
          _block(height: 26, width: 180),
          const SizedBox(height: AppSizes.sm),
          _block(height: 20),
          const SizedBox(height: AppSizes.sm),
          _block(height: 20, width: 260),
          const SizedBox(height: AppSizes.lg),
          Row(
            children: [
              Expanded(child: _block(height: 78)),
              const SizedBox(width: AppSizes.sm),
              Expanded(child: _block(height: 78)),
            ],
          ),
          const SizedBox(height: AppSizes.sm),
          Row(
            children: [
              Expanded(child: _block(height: 78)),
              const SizedBox(width: AppSizes.sm),
              Expanded(child: _block(height: 78)),
            ],
          ),
          const SizedBox(height: AppSizes.lg),
          _block(height: 76),
          const SizedBox(height: AppSizes.md),
          _block(height: 76),
          const SizedBox(height: AppSizes.md),
          _block(height: 76),
        ],
      ),
    );
  }

  Widget _block({required double height, double? width}) {
    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: AppColors.bgSecondary,
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      ),
    );
  }
}
