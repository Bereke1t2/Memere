"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { 
  Menu, 
  X, 
  ArrowRight, 
  Smartphone, 
  Sparkles, 
  Send 
} from "lucide-react";

export function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [activeSection, setActiveSection] = useState("hero");

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20);

      const sections = ["hero", "curriculum", "how-it-works", "features", "mock-exams", "testimonials", "faq"];
      const current = sections.find((sec) => {
        const el = document.getElementById(sec);
        if (el) {
          const rect = el.getBoundingClientRect();
          return rect.top <= 120 && rect.bottom >= 120;
        }
        return false;
      });
      if (current) setActiveSection(current);
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const navLinks = [
    { label: "Home", href: "#hero", id: "hero" },
    { label: "What we do", href: "#curriculum", id: "curriculum" },
    { label: "Our Approach", href: "#how-it-works", id: "how-it-works" },
    { label: "Mock Exams", href: "#mock-exams", id: "mock-exams" },
    { label: "Testimonials", href: "#testimonials", id: "testimonials" },
    { label: "FAQ", href: "#faq", id: "faq" },
  ];

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-white/95 backdrop-blur-md shadow-[0_4px_20px_rgba(30,43,88,0.06)] py-3.5 border-b border-[#F1F5F9]"
          : "bg-transparent py-5"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          
          {/* Logo with Brand Icon */}
          <Link href="/" className="flex items-center gap-3 group">
            <Image
              src="/logo.png"
              alt="Mirkuz Logo"
              width={38}
              height={38}
              priority
              className="w-9 h-9 rounded-xl shadow-md shadow-[#1E2B58]/15 group-hover:scale-105 transition-transform object-contain"
            />
            <div className="flex items-baseline gap-1.5">
              <span className="text-2xl font-extrabold tracking-tight text-[#1E2B58]">
                Mirkuz
              </span>
              <span className="text-[11px] font-bold text-[#00B894] bg-[#E8F8F5] px-2 py-0.5 rounded-full uppercase tracking-wider">
                ምርኩዝ
              </span>
            </div>
          </Link>

          {/* Desktop Nav Items with Active Teal Dot */}
          <nav className="hidden md:flex items-center gap-9">
            {navLinks.map((link) => {
              const isActive = activeSection === link.id;
              return (
                <a
                  key={link.label}
                  href={link.href}
                  className={`text-sm font-semibold transition-colors relative py-1 ${
                    isActive ? "text-[#1E2B58]" : "text-[#64748B] hover:text-[#1E2B58]"
                  }`}
                >
                  <span>{link.label}</span>
                  {/* Active Teal Indicator Dot (as in video) */}
                  {isActive && (
                    <span className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-1.5 h-1.5 rounded-full bg-[#00B894]" />
                  )}
                </a>
              );
            })}
          </nav>

          {/* Right Action CTAs */}
          <div className="hidden md:flex items-center gap-4">
            <a
              href="https://t.me/mirkuz_exam"
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs font-bold text-[#1E2B58] hover:text-[#00B894] px-3 py-2 transition-colors flex items-center gap-1.5"
            >
              <Send className="w-3.5 h-3.5 text-[#00B894]" />
              <span>Community</span>
            </a>
            
            {/* Deep Navy Button with Teal hover accent (Matching video Frame 002) */}
            <a
              href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#1E2B58] hover:bg-[#162145] text-white text-xs sm:text-sm font-bold shadow-md shadow-[#1E2B58]/20 transition-all hover:shadow-lg hover:-translate-y-0.5 active:translate-y-0"
            >
              <Smartphone className="w-4 h-4 text-[#00B894]" />
              <span>Get started</span>
              <ArrowRight className="w-3.5 h-3.5 text-white/80" />
            </a>
          </div>

          {/* Mobile Menu Hamburger Button */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-xl text-[#1E2B58] hover:bg-slate-100 transition-colors"
            aria-label="Toggle navigation menu"
          >
            {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>

        </div>
      </div>

      {/* Mobile Drawer Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-white/98 backdrop-blur-xl border-b border-[#F1F5F9] shadow-2xl px-6 py-6 space-y-4 animate-in slide-in-from-top-2 duration-200">
          <nav className="flex flex-col space-y-3">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                onClick={() => setMobileMenuOpen(false)}
                className="text-base font-semibold text-[#1E2B58] hover:text-[#00B894] py-1.5 transition-colors flex items-center justify-between"
              >
                <span>{link.label}</span>
                <ArrowRight className="w-4 h-4 text-[#94A3B8]" />
              </a>
            ))}
          </nav>
          
          <div className="pt-4 border-t border-[#F1F5F9] flex flex-col gap-3">
            <a
              href="https://t.me/mirkuz_exam"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full text-center py-2.5 text-xs font-bold text-[#1E2B58] border border-[#CBD5E1] rounded-xl hover:bg-slate-50 transition-colors flex items-center justify-center gap-2"
            >
              <Send className="w-3.5 h-3.5 text-[#00B894]" />
              <span>Join Telegram Channel</span>
            </a>
            <a
              href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full text-center py-3 text-xs font-bold text-white bg-[#1E2B58] hover:bg-[#162145] rounded-xl shadow-md shadow-[#1E2B58]/20 transition-all flex items-center justify-center gap-2"
            >
              <Smartphone className="w-4 h-4 text-[#00B894]" />
              <span>Download App on Google Play</span>
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
