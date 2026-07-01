import 'package:flutter/material.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_text_styles.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../domain/entities/payment_provider_entity.dart';

/// Bottom sheet for choosing a payment provider. Never collects card details —
/// the provider's own hosted checkout handles that. Returns the chosen provider,
/// or null if the user cancels.
class PaymentProviderSheet extends StatefulWidget {
  const PaymentProviderSheet({super.key, this.amountLabel});

  final String? amountLabel;

  static Future<PaymentProvider?> show(
    BuildContext context, {
    String? amountLabel,
  }) {
    return showModalBottomSheet<PaymentProvider>(
      context: context,
      backgroundColor: AppColors.bgSecondary,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSizes.radiusLg)),
      ),
      builder: (_) => PaymentProviderSheet(amountLabel: amountLabel),
    );
  }

  @override
  State<PaymentProviderSheet> createState() => _PaymentProviderSheetState();
}

class _PaymentProviderSheetState extends State<PaymentProviderSheet> {
  late PaymentProvider _selected = PaymentProviders.visible().first;

  @override
  Widget build(BuildContext context) {
    final providers = PaymentProviders.visible();

    return SafeArea(
      top: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(
          AppSizes.md,
          AppSizes.sm,
          AppSizes.md,
          AppSizes.md,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Center(
              child: Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.only(bottom: AppSizes.md),
                decoration: BoxDecoration(
                  color: AppColors.border,
                  borderRadius: BorderRadius.circular(AppSizes.radiusFull),
                ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Choose payment method',
                    style: AppTextStyles.headlineSmall),
                if (widget.amountLabel != null)
                  Text(
                    widget.amountLabel!,
                    style: AppTextStyles.titleMedium.copyWith(
                      color: AppColors.accentSecondary,
                    ),
                  ),
              ],
            ),
            const SizedBox(height: AppSizes.md),
            ...providers.map(
              (provider) => Padding(
                padding: const EdgeInsets.only(bottom: AppSizes.sm),
                child: _ProviderOption(
                  provider: provider,
                  selected: provider == _selected,
                  onTap: () => setState(() => _selected = provider),
                ),
              ),
            ),
            const SizedBox(height: AppSizes.sm),
            Row(
              children: [
                Expanded(
                  child: AppButton(
                    label: 'Cancel',
                    variant: AppButtonVariant.outline,
                    onPressed: () => Navigator.of(context).pop(),
                  ),
                ),
                const SizedBox(width: AppSizes.md),
                Expanded(
                  child: AppButton(
                    label: 'Continue',
                    onPressed: () => Navigator.of(context).pop(_selected),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _ProviderOption extends StatelessWidget {
  const _ProviderOption({
    required this.provider,
    required this.selected,
    required this.onTap,
  });

  final PaymentProvider provider;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      borderRadius: BorderRadius.circular(AppSizes.radiusMd),
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(AppSizes.md),
        decoration: BoxDecoration(
          color: AppColors.bgTertiary,
          borderRadius: BorderRadius.circular(AppSizes.radiusMd),
          border: Border.all(
            color: selected ? AppColors.accentPrimary : AppColors.border,
            width: selected ? 1.5 : 1,
          ),
        ),
        child: Row(
          children: [
            Container(
              width: AppSizes.avatarSm,
              height: AppSizes.avatarSm,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: AppColors.accentGlow,
                borderRadius: BorderRadius.circular(AppSizes.radiusSm),
              ),
              child: Text(
                provider.label.substring(0, 1),
                style: AppTextStyles.titleMedium
                    .copyWith(color: AppColors.accentPrimary),
              ),
            ),
            const SizedBox(width: AppSizes.md),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(provider.label, style: AppTextStyles.titleMedium),
                  const SizedBox(height: AppSizes.xs),
                  Text(provider.description, style: AppTextStyles.bodySmall),
                ],
              ),
            ),
            Icon(
              selected
                  ? Icons.radio_button_checked_rounded
                  : Icons.radio_button_unchecked_rounded,
              color: selected ? AppColors.accentPrimary : AppColors.textDisabled,
              size: AppSizes.iconMd,
            ),
          ],
        ),
      ),
    );
  }
}
