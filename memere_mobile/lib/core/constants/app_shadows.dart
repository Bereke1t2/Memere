import 'package:flutter/material.dart';

/// Layered shadow tokens for the app's dark surfaces.
///
/// Dark UI needs a tight contact shadow plus a wider ambient shadow. A subtle
/// green cast keeps elevated surfaces connected to the app palette.
abstract class AppShadows {
  static const List<BoxShadow> sm = [
    BoxShadow(
      color: Color(0x3D000000),
      blurRadius: 2,
      offset: Offset(0, 1),
    ),
    BoxShadow(
      color: Color(0x2604110C),
      blurRadius: 10,
      offset: Offset(0, 4),
    ),
  ];

  static const List<BoxShadow> md = [
    BoxShadow(
      color: Color(0x47000000),
      blurRadius: 4,
      offset: Offset(0, 2),
    ),
    BoxShadow(
      color: Color(0x3004110C),
      blurRadius: 18,
      offset: Offset(0, 8),
    ),
  ];

  static const List<BoxShadow> lg = [
    BoxShadow(
      color: Color(0x59000000),
      blurRadius: 8,
      offset: Offset(0, 4),
    ),
    BoxShadow(
      color: Color(0x3804110C),
      blurRadius: 34,
      offset: Offset(0, 16),
    ),
  ];

  static const List<BoxShadow> accentGlow = [
    BoxShadow(
      color: Color(0x3335B87E),
      blurRadius: 22,
      offset: Offset(0, 8),
    ),
  ];

  static const List<BoxShadow> focus = [
    BoxShadow(
      color: Color(0x2435B87E),
      blurRadius: 18,
      spreadRadius: 1,
      offset: Offset(0, 4),
    ),
  ];
}
