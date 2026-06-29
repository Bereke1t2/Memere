import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

import '../../../../core/constants/app_colors.dart';

class HlsVideoPlayer extends StatelessWidget {
  const HlsVideoPlayer({
    super.key,
    required this.controller,
    required this.chewieController,
  });

  final VideoPlayerController controller;
  final ChewieController chewieController;

  @override
  Widget build(BuildContext context) {
    if (!controller.value.isInitialized) {
      return const AspectRatio(
        aspectRatio: 16 / 9,
        child: ColoredBox(
          color: Colors.black,
          child: Center(
            child: CircularProgressIndicator(
              color: AppColors.accentPrimary,
              strokeWidth: 2,
            ),
          ),
        ),
      );
    }

    return AspectRatio(
      aspectRatio: controller.value.aspectRatio == 0
          ? 16 / 9
          : controller.value.aspectRatio,
      child: Chewie(controller: chewieController),
    );
  }
}
