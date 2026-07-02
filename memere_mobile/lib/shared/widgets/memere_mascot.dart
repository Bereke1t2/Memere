import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';

class MemereMascot extends StatefulWidget {
  const MemereMascot({
    super.key,
    this.size = const Size(280, 252),
    this.showBackdrop = true,
  });

  final Size size;
  final bool showBackdrop;

  @override
  State<MemereMascot> createState() => _MemereMascotState();
}

class _MemereMascotState extends State<MemereMascot>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 3200),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.size.width,
      height: widget.size.height,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, child) {
          final sway = math.sin(_controller.value * math.pi * 2);
          return Transform.translate(
            offset: Offset(0, sway * 2),
            child: Transform.rotate(
              angle: sway * 0.008,
              child: child,
            ),
          );
        },
        child: RepaintBoundary(
          child: CustomPaint(
            painter: _MemereMascotPainter(
              phase: _controller,
              showBackdrop: widget.showBackdrop,
            ),
          ),
        ),
      ),
    );
  }
}

class _MemereMascotPainter extends CustomPainter {
  _MemereMascotPainter({
    required this.phase,
    required this.showBackdrop,
  }) : super(repaint: phase);

  final Animation<double> phase;
  final bool showBackdrop;

  @override
  void paint(Canvas canvas, Size size) {
    const base = Size(320, 280);
    final scale = math.min(size.width / base.width, size.height / base.height);
    final dx = (size.width - base.width * scale) / 2;
    final dy = (size.height - base.height * scale) / 2;

    canvas.save();
    canvas.translate(dx, dy);
    canvas.scale(scale);

    final value = phase.value;
    final breathe = math.sin(value * math.pi * 2);
    final blink = ((value * 3) % 1) > 0.92 ? 0.12 : 1.0;

    if (showBackdrop) _drawBackdrop(canvas);
    _drawDesk(canvas);
    _drawTeacher(canvas, breathe, blink);
    _drawLaptopAndBook(canvas);

    canvas.restore();
  }

  void _drawBackdrop(Canvas canvas) {
    final backdropPaint = Paint()
      ..shader = AppColors.headerGradient.createShader(
        const Rect.fromLTWH(26, 12, 268, 232),
      );
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(26, 12, 268, 232),
        const Radius.circular(42),
      ),
      backdropPaint,
    );

    final glowPaint = Paint()..color = AppColors.accentPrimary.withAlpha(26);
    canvas.drawCircle(const Offset(95, 74), 46, glowPaint);
    canvas.drawCircle(const Offset(245, 72), 28, glowPaint);

    final trimPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 4
      ..strokeCap = StrokeCap.round;
    for (final item in const [
      (AppColors.coral, 63.0),
      (AppColors.warning, 74.0),
      (AppColors.success, 85.0),
    ]) {
      trimPaint.color = item.$1.withAlpha(160);
      canvas.drawLine(
        Offset(210, item.$2),
        Offset(278, item.$2 + 18),
        trimPaint,
      );
    }
  }

  void _drawDesk(Canvas canvas) {
    final shadow = Paint()..color = const Color(0x263F168B);
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(48, 230, 224, 18),
        const Radius.circular(16),
      ),
      shadow,
    );

    final deskPaint = Paint()..color = AppColors.bgQuaternary;
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(42, 210, 236, 34),
        const Radius.circular(17),
      ),
      deskPaint,
    );
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(66, 238, 18, 32),
        const Radius.circular(9),
      ),
      deskPaint,
    );
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(236, 238, 18, 32),
        const Radius.circular(9),
      ),
      deskPaint,
    );
  }

  void _drawTeacher(Canvas canvas, double breathe, double blink) {
    final skin = Paint()..color = const Color(0xFF84614B);
    final skinLight = Paint()..color = const Color(0xFF9A745C);
    final hair = Paint()..color = const Color(0xFFCFC8BA);
    final hairShadow = Paint()..color = const Color(0xFF9F9A8E);
    final line = Paint()
      ..color = AppColors.bgPrimary
      ..style = PaintingStyle.stroke
      ..strokeWidth = 3
      ..strokeCap = StrokeCap.round;

    final bodyLift = breathe * 2;
    final bodyPath = Path()
      ..moveTo(105, 204 - bodyLift)
      ..quadraticBezierTo(160, 160 - bodyLift, 215, 204 - bodyLift)
      ..lineTo(225, 220)
      ..lineTo(95, 220)
      ..close();
    canvas.drawPath(bodyPath, Paint()..color = AppColors.textSecondary);

    final shawl = Path()
      ..moveTo(111, 203 - bodyLift)
      ..quadraticBezierTo(160, 166 - bodyLift, 209, 203 - bodyLift)
      ..lineTo(198, 220)
      ..quadraticBezierTo(160, 203 - bodyLift, 122, 220)
      ..close();
    canvas.drawPath(shawl, Paint()..color = AppColors.accentPrimaryDeep);

    final trim = Paint()
      ..strokeWidth = 3
      ..strokeCap = StrokeCap.round
      ..style = PaintingStyle.stroke;
    for (final item in const [
      (AppColors.coral, -8.0),
      (AppColors.warning, 0.0),
      (AppColors.success, 8.0),
    ]) {
      trim.color = item.$1;
      canvas.drawLine(
        Offset(126 + item.$2, 191 - bodyLift),
        Offset(141 + item.$2, 218),
        trim,
      );
    }

    canvas.drawOval(
      const Rect.fromLTWH(95, 184, 52, 28),
      skinLight,
    );
    canvas.drawOval(
      const Rect.fromLTWH(173, 184, 52, 28),
      skinLight,
    );

    canvas.drawOval(
      const Rect.fromLTWH(103, 68, 114, 118),
      hairShadow,
    );
    canvas.drawOval(
      const Rect.fromLTWH(97, 56, 126, 124),
      hair,
    );

    canvas.drawOval(
      const Rect.fromLTWH(111, 76, 98, 98),
      skin,
    );
    canvas.drawOval(
      const Rect.fromLTWH(103, 112, 18, 28),
      skin,
    );
    canvas.drawOval(
      const Rect.fromLTWH(199, 112, 18, 28),
      skin,
    );

    canvas.drawArc(
      const Rect.fromLTWH(111, 55, 98, 74),
      math.pi,
      math.pi,
      false,
      Paint()
        ..color = hair.color
        ..style = PaintingStyle.stroke
        ..strokeWidth = 22
        ..strokeCap = StrokeCap.round,
    );

    final glasses = Paint()
      ..color = AppColors.bgPrimary
      ..style = PaintingStyle.stroke
      ..strokeWidth = 3;
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(122, 110, 32, 23),
        const Radius.circular(10),
      ),
      glasses,
    );
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(166, 110, 32, 23),
        const Radius.circular(10),
      ),
      glasses,
    );
    canvas.drawLine(const Offset(154, 121), const Offset(166, 121), glasses);

    final eyePaint = Paint()
      ..color = AppColors.bgPrimary
      ..strokeWidth = 2
      ..strokeCap = StrokeCap.round;
    canvas.drawLine(
      const Offset(132, 121),
      Offset(144, 121 + (1 - blink) * 4),
      eyePaint,
    );
    canvas.drawLine(
      Offset(176, 121 + (1 - blink) * 4),
      const Offset(188, 121),
      eyePaint,
    );

    canvas.drawLine(const Offset(160, 124), const Offset(156, 140), line);
    canvas.drawArc(
      const Rect.fromLTWH(142, 141, 36, 18),
      0.18,
      math.pi - 0.36,
      false,
      line,
    );
  }

  void _drawLaptopAndBook(Canvas canvas) {
    final laptop = Paint()..color = AppColors.bgPrimary;
    final screen = RRect.fromRectAndRadius(
      const Rect.fromLTWH(52, 166, 82, 54),
      const Radius.circular(10),
    );
    canvas.drawRRect(screen, laptop);
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(64, 176, 58, 32),
        const Radius.circular(6),
      ),
      Paint()..color = AppColors.accentPrimary.withAlpha(180),
    );
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        const Rect.fromLTWH(44, 218, 100, 9),
        const Radius.circular(5),
      ),
      Paint()..color = AppColors.bgQuaternary,
    );

    final bookLeft = Path()
      ..moveTo(187, 184)
      ..quadraticBezierTo(216, 173, 245, 184)
      ..lineTo(245, 221)
      ..quadraticBezierTo(216, 210, 187, 221)
      ..close();
    final bookRight = Path()
      ..moveTo(245, 184)
      ..quadraticBezierTo(265, 177, 282, 188)
      ..lineTo(282, 223)
      ..quadraticBezierTo(265, 215, 245, 221)
      ..close();
    canvas.drawPath(bookLeft, Paint()..color = const Color(0xFFC5B996));
    canvas.drawPath(bookRight, Paint()..color = const Color(0xFFB5A87F));
    canvas.drawLine(
      const Offset(245, 184),
      const Offset(245, 221),
      Paint()
        ..color = AppColors.accentSecondary
        ..strokeWidth = 2,
    );
  }

  @override
  bool shouldRepaint(covariant _MemereMascotPainter oldDelegate) {
    return oldDelegate.showBackdrop != showBackdrop ||
        oldDelegate.phase != phase;
  }
}
