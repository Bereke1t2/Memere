import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

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
      return Container(
        width: double.infinity,
        decoration: BoxDecoration(
          color: Colors.black,
          borderRadius: BorderRadius.circular(12),
        ),
        child: const AspectRatio(
          aspectRatio: 16 / 9,
          child: Center(
            child: CircularProgressIndicator(
              color: Color(0xFFFF5252),
              strokeWidth: 2,
            ),
          ),
        ),
      );
    }

    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        color: Colors.black,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withAlpha(120),
            blurRadius: 16,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      clipBehavior: Clip.antiAlias,
      child: AspectRatio(
        aspectRatio: 16 / 9,
        child: Chewie(controller: chewieController),
      ),
    );
  }
}
