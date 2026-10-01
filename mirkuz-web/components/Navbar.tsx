"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { 
  BookOpen, 
  Download, 
  Menu, 
  X, 
  Send, 
  Sparkles, 
  ChevronRight,
  GraduationCap
} from "lucide-react";

export function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      if (window.scrollY > 20) {
        setIsScrolled(true);
      } else {
        setIsScrolled(false);
      }
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-[#050505]/90 backdrop-blur-xl border-b border-[#1a1a20] shadow-2xl shadow-black/40 py-3"
          : "bg-transparent py-5"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          {/* Brand Logo */}
          <Link href="/" className="flex items-center gap-3 group">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#10b981] to-[#059669] flex items-center justify-center shadow-lg shadow-[#10b981]/20 group-hover:scale-105 transition-transform">
              <GraduationCap className="w-5 h-5 text-white" />
            </div>
            <div className="flex flex-col">
              <div className="flex items-center gap-2">
                <span className="text-xl font-bold tracking-tight text-white group-hover:text-[#10b981] transition-colors">
                  Mirkuz
                </span>
                <span className="px-1.5 py-0.5 text-[10px] font-semibold bg-[#10b981]/15 text-[#10b981] rounded border border-[#10b981]/30">
                  ምርኩዝ
                </span>
              </div>
              <span className="text-[11px] text-[#94a3b8] font-medium tracking-wide">
                Grade 12 Entrance Prep
              </span>
            </div>
          </Link>

          {/* Desktop Nav Links */}
          <nav className="hidden md:flex items-center gap-8">
            <a
              href="#curriculum"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors"
            >
              Curriculum & Streams
            </a>
            <a
              href="#features"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors"
            >
              Features
            </a>
            <a
              href="#exam-demo"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors flex items-center gap-1.5"
            >
              <span className="w-2 h-2 rounded-full bg-[#10b981] animate-pulse"></span>
              Live Mock Demo
            </a>
            <a
              href="#offline"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors"
            >
              Offline Study
            </a>
            <a
              href="#pricing"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors"
            >
              Pricing & Telebirr
            </a>
            <a
              href="#faq"
              className="text-sm font-medium text-[#94a3b8] hover:text-white transition-colors"
            >
              FAQ
            </a>
          </nav>

          {/* Desktop Action CTAs */}
          <div className="hidden lg:flex items-center gap-3">
            <a
              href="https://t.me/mirkuz_exam"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-3.5 py-2 rounded-lg text-xs font-semibold bg-[#18181e] text-[#94a3b8] hover:text-white hover:bg-[#25252e] border border-[#2c2c35] transition-all"
            >
              <Send className="w-3.5 h-3.5 text-[#38bdf8]" />
              <span>Join Telegram</span>
            </a>

            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold bg-gradient-to-r from-[#10b981] to-[#059669] text-white hover:opacity-95 shadow-md shadow-[#10b981]/25 hover:shadow-lg hover:shadow-[#10b981]/40 transition-all active:scale-[0.98]"
            >
              <Download className="w-3.5 h-3.5" />
              <span>Get on Google Play</span>
            </a>
          </div>

          {/* Mobile Hamburger Button */}
          <div className="flex md:hidden items-center gap-2">
            <a
              href="https://t.me/mirkuz_exam"
              target="_blank"
              rel="noopener noreferrer"
              className="p-2 rounded-lg bg-[#18181e] text-[#38bdf8] border border-[#2c2c35]"
              aria-label="Telegram"
            >
              <Send className="w-4 h-4" />
            </a>

            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="p-2 rounded-lg bg-[#111115] text-[#94a3b8] hover:text-white border border-[#1a1a20]"
              aria-label="Toggle navigation menu"
            >
              {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Drawer */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-[#0a0a0d] border-b border-[#1a1a20] px-4 pt-4 pb-6 mt-3 space-y-4 animate-in slide-in-from-top-4 duration-200">
          <nav className="flex flex-col space-y-3">
            <a
              href="#curriculum"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#94a3b8] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span>Curriculum & Streams</span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
            <a
              href="#features"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#94a3b8] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span>Key Features</span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
            <a
              href="#exam-demo"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#10b981] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-[#10b981]"></span>
                Live Mock Demo
              </span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
            <a
              href="#offline"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#94a3b8] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span>100% Offline Study</span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
            <a
              href="#pricing"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#94a3b8] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span>Pricing & Telebirr</span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
            <a
              href="#faq"
              onClick={() => setMobileMenuOpen(false)}
              className="text-sm font-medium text-[#94a3b8] hover:text-white py-1.5 flex items-center justify-between"
            >
              <span>FAQ</span>
              <ChevronRight className="w-4 h-4 text-[#71717a]" />
            </a>
          </nav>

          <div className="pt-3 border-t border-[#1a1a20] flex flex-col gap-2.5">
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full flex items-center justify-center gap-2 py-2.5 rounded-lg text-sm font-semibold bg-gradient-to-r from-[#10b981] to-[#059669] text-white shadow-lg shadow-[#10b981]/25"
            >
              <Download className="w-4 h-4" />
              <span>Download on Google Play</span>
            </a>

            <a
              href="https://t.me/mirkuz_exam"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full flex items-center justify-center gap-2 py-2.5 rounded-lg text-sm font-semibold bg-[#18181e] text-white border border-[#2c2c35]"
            >
              <Send className="w-4 h-4 text-[#38bdf8]" />
              <span>Join Telegram Community</span>
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
