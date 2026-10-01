"use client";

import React from "react";
import Link from "next/link";
import { GraduationCap, Send, Download, Mail, Phone, ExternalLink, ShieldCheck } from "lucide-react";

export function Footer() {
  return (
    <footer className="bg-[#050505] border-t border-[#1a1a20] pt-16 pb-12 text-[#94a3b8] text-xs">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8 pb-12 border-b border-[#1a1a20]">
          
          {/* Brand Col (Span 2) */}
          <div className="lg:col-span-2 space-y-4">
            <Link href="/" className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#10b981] to-[#059669] flex items-center justify-center text-white font-bold">
                <GraduationCap className="w-5 h-5" />
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <span className="text-xl font-bold tracking-tight text-white">Mirkuz</span>
                  <span className="px-1.5 py-0.5 text-[10px] font-semibold bg-[#10b981]/15 text-[#10b981] rounded border border-[#10b981]/30">
                    ምርኩዝ
                  </span>
                </div>
                <span className="text-[11px] text-[#71717a]">
                  Grade 12 University Entrance Learning Platform
                </span>
              </div>
            </Link>

            <p className="text-xs text-[#94a3b8] max-w-sm leading-relaxed">
              Empowering Ethiopian Grade 12 students with top-tier video lessons, downloadable chapter summaries, timed national mock entrance exams, and 100% offline study capability.
            </p>

            <div className="flex items-center gap-3 pt-2">
              <a
                href="https://t.me/mirkuz_exam"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-lg bg-[#18181e] border border-[#2c2c35] flex items-center justify-center text-[#38bdf8] hover:bg-[#25252e] transition-colors"
                aria-label="Telegram"
              >
                <Send className="w-4 h-4" />
              </a>

              <a
                href="https://play.google.com/store/apps"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-lg bg-[#18181e] border border-[#2c2c35] flex items-center justify-center text-[#10b981] hover:bg-[#25252e] transition-colors"
                aria-label="Google Play"
              >
                <Download className="w-4 h-4" />
              </a>
            </div>
          </div>

          {/* Col 2: Natural Science */}
          <div className="space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-white">
              🌿 Natural Science
            </h4>
            <ul className="space-y-2">
              <li><a href="#curriculum" className="hover:text-white transition-colors">Advanced Mathematics</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Physics</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Chemistry</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Biology</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">English (Natural)</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Civics & Aptitude</a></li>
            </ul>
          </div>

          {/* Col 3: Social Science */}
          <div className="space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-white">
              🏛️ Social Science
            </h4>
            <ul className="space-y-2">
              <li><a href="#curriculum" className="hover:text-white transition-colors">General Mathematics</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">History of Ethiopia & World</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Geography</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Economics</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">English (Social)</a></li>
              <li><a href="#curriculum" className="hover:text-white transition-colors">Civics & Aptitude</a></li>
            </ul>
          </div>

          {/* Col 4: Platform & Support */}
          <div className="space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-white">
              Platform & Access
            </h4>
            <ul className="space-y-2">
              <li><a href="#exam-demo" className="hover:text-white transition-colors">Live Mock Demo</a></li>
              <li><a href="#offline" className="hover:text-white transition-colors">100% Offline Mode</a></li>
              <li><a href="#pricing" className="hover:text-white transition-colors">Telebirr & Pricing</a></li>
              <li><a href="#faq" className="hover:text-white transition-colors">Frequently Asked Questions</a></li>
              <li>
                <a
                  href="http://localhost:3000"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:text-white transition-colors inline-flex items-center gap-1 text-[#10b981]"
                >
                  <span>Teacher & Admin Portal</span>
                  <ExternalLink className="w-3 h-3" />
                </a>
              </li>
            </ul>
          </div>

        </div>

        {/* Bottom Bar */}
        <div className="pt-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-[11px] text-[#71717a]">
          <p>© {new Date().getFullYear()} Mirkuz (ምርኩዝ) Education Platform. All rights reserved.</p>
          <div className="flex items-center gap-6">
            <span>Made with ❤️ for Ethiopian Students</span>
            <span>Addis Ababa, Ethiopia</span>
          </div>
        </div>

      </div>
    </footer>
  );
}
