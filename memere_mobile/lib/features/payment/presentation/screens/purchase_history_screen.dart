import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/router/app_router.dart';
import '../providers/purchase_history_provider.dart';
import '../widgets/enrollment_tile.dart';
import '../widgets/payment_empty_state.dart';
import '../widgets/payment_history_tile.dart';

class PurchaseHistoryScreen extends ConsumerWidget {
  const PurchaseHistoryScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        backgroundColor: AppColors.bgPrimary,
        appBar: AppBar(
          title: const Text('Purchases'),
          bottom: const TabBar(
            indicatorColor: AppColors.accentPrimary,
            labelColor: AppColors.textPrimary,
            unselectedLabelColor: AppColors.textSecondary,
            tabs: [
              Tab(text: 'Payments'),
              Tab(text: 'Enrollments'),
            ],
          ),
        ),
        body: SafeArea(
          top: false,
          child: TabBarView(
            children: [
              _PaymentsTab(),
              _EnrollmentsTab(),
            ],
          ),
        ),
      ),
    );
  }
}

class _PaymentsTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(paymentHistoryProvider);

    return async.when(
      loading: () => const _ListSkeleton(),
      error: (error, _) => _ErrorState(
        message: error is Failure ? error.message : 'Could not load payments.',
        onRetry: () => ref.invalidate(paymentHistoryProvider),
      ),
      data: (payments) {
        if (payments.isEmpty) {
          return const PaymentEmptyState(
            icon: Icons.receipt_long_outlined,
            title: 'No payments yet',
            body: 'Your course and subscription payments will appear here.',
          );
        }
        return RefreshIndicator(
          color: AppColors.accentPrimary,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(paymentHistoryProvider);
            await ref.read(paymentHistoryProvider.future);
          },
          child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(AppSizes.md),
            itemCount: payments.length,
            separatorBuilder: (_, __) => const SizedBox(height: AppSizes.sm),
            itemBuilder: (_, index) {
              final payment = payments[index];
              return PaymentHistoryTile(
                payment: payment,
                onTap: payment.courseId != null && payment.courseId!.isNotEmpty
                    ? () => context.push(
                          AppRoutes.courseDetailPath(payment.courseId!),
                        )
                    : null,
              );
            },
          ),
        );
      },
    );
  }
}

class _EnrollmentsTab extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(enrollmentListProvider);

    return async.when(
      loading: () => const _ListSkeleton(),
      error: (error, _) => _ErrorState(
        message:
            error is Failure ? error.message : 'Could not load enrollments.',
        onRetry: () => ref.invalidate(enrollmentListProvider),
      ),
      data: (enrollments) {
        if (enrollments.isEmpty) {
          return const PaymentEmptyState(
            icon: Icons.school_outlined,
            title: 'No enrollments yet',
            body: 'Courses you enroll in will appear here.',
          );
        }
        return RefreshIndicator(
          color: AppColors.accentPrimary,
          backgroundColor: AppColors.bgSecondary,
          onRefresh: () async {
            ref.invalidate(enrollmentListProvider);
            await ref.read(enrollmentListProvider.future);
          },
          child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(AppSizes.md),
            itemCount: enrollments.length,
            separatorBuilder: (_, __) => const SizedBox(height: AppSizes.sm),
            itemBuilder: (_, index) {
              final enrollment = enrollments[index];
              return EnrollmentTile(
                enrollment: enrollment,
                onTap: () => context.push(
                  AppRoutes.courseDetailPath(enrollment.courseId),
                ),
              );
            },
          ),
        );
      },
    );
  }
}

class _ErrorState extends StatelessWidget {
  const _ErrorState({required this.message, required this.onRetry});

  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return PaymentEmptyState(
      icon: Icons.error_outline_rounded,
      title: 'Something went wrong',
      body: message,
      buttonLabel: 'Retry',
      onPressed: onRetry,
      iconColor: AppColors.error,
    );
  }
}

class _ListSkeleton extends StatelessWidget {
  const _ListSkeleton();

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
      padding: const EdgeInsets.all(AppSizes.md),
      itemCount: 6,
      separatorBuilder: (_, __) => const SizedBox(height: AppSizes.sm),
      itemBuilder: (_, __) => Container(
        height: 72,
        decoration: BoxDecoration(
          color: AppColors.bgSecondary,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
      ),
    );
  }
}
