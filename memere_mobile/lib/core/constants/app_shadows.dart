import 'package:flutter/material.dart';

/// Subtle elevation tokens for dark learning surfaces.
abstract class AppShadows {
  static const List<BoxShadow> sm = [
    BoxShadow(
      color: Color(0x40000000),
      blurRadius: 10,
      offset: Offset(0, 4),
    ),
  ];

  static const List<BoxShadow> md = [
    BoxShadow(
      color: Color(0x52000000),
      blurRadius: 18,
      offset: Offset(0, 8),
    ),
  ];

  static const List<BoxShadow> lg = [
    BoxShadow(
      color: Color(0x66000000),
      blurRadius: 28,
      offset: Offset(0, 10),
    ),
  ];

  static const List<BoxShadow> accentGlow = [
    BoxShadow(
      color: Color(0x22FFFFFF),
      blurRadius: 18,
      offset: Offset(0, 8),
    ),
  ];

  static const List<BoxShadow> focus = [
    BoxShadow(
      color: Color(0x33FFFFFF),
      blurRadius: 16,
      spreadRadius: 1,
      offset: Offset(0, 6),
    ),
  ];
}
