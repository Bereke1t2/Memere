package com.example.memere_mobile

import android.os.Bundle
import android.view.WindowManager
import io.flutter.embedding.android.FlutterActivity

class MainActivity : FlutterActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        // Prevent screenshots, screen recording, and app switcher preview capture
        window.addFlags(WindowManager.LayoutParams.FLAG_SECURE)
    }
}

