"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { 
  Menu, 
  X, 
  ArrowRight,
  Sparkles,
  Smartphone,
  GraduationCap
} from "lucide-react";

export function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 15);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const navLinks = [
    { label: "Home", href: "#hero" },
    { label: "Service", href: "#services" },
    { label: "Feature", href: "#features" },
    { label: "Product", href: "#product" },
    { label: "Testimonial", href: "#testimonial" },
    { label: "FAQ", href: "#faq" },
  ];

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-white/95 backdrop-blur-md shadow-[0_2px_10px_rgba(0,0,0,0.06)] py-3.5 border-b border-[#E8ECF2]"
          : "bg-[#F5F7FA] py-5"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2.5 group">
            <div className="w-8 h-8 rounded-lg bg-[#4CAF4F] flex items-center justify-center text-white shadow-sm group-hover:scale-105 transition-transform">
              <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
                <path d="M12 2L2 7l10 5 10-5-10-5zm0 9.2L4.5 7.4 12 3.6l7.5 3.8L12 11.2zm0 2.8L2 9v6l10 5 10-5V9l-10 5zm0 2.2l-7.5-3.8V11l7.5 3.8 7.5-3.8v2.6L12 16.2z"/>
              </svg>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-bold tracking-tight text-[#263238]">
                Mirkuz
              </span>
              <span className="text-xs font-semibold text-[#4CAF4F] hidden sm:inline">
                Nexcent
              </span>
            </div>
          </Link>

          {/* Desktop Nav Items */}
          <nav className="hidden md:flex items-center gap-8">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                className="text-sm font-medium text-[#4D4D4D] hover:text-[#4CAF4F] transition-colors"
              >
                {link.label}
              </a>
            ))}
          </nav>

          {/* Right Action CTAs */}
          <div className="hidden md:flex items-center gap-4">
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm font-medium text-[#4CAF4F] hover:text-[#388E3C] px-3 py-2 transition-colors"
            >
              Login
            </a>
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-md bg-[#4CAF4F] text-white text-sm font-medium hover:bg-[#388E3C] shadow-sm transition-all hover:shadow hover:-translate-y-0.5 active:translate-y-0"
            >
              <span>Sign up</span>
              <ArrowRight className="w-4 h-4" />
            </a>
          </div>

          {/* Mobile Menu Hamburger Button */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-lg text-[#4D4D4D] hover:bg-black/5 transition-colors"
            aria-label="Toggle navigation menu"
          >
            {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>

        </div>
      </div>

      {/* Mobile Drawer Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-white border-b border-[#E8ECF2] shadow-xl px-6 py-6 space-y-4 animate-in slide-in-from-top-2 duration-200">
          <nav className="flex flex-col space-y-3">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                onClick={() => setMobileMenuOpen(false)}
                className="text-base font-medium text-[#4D4D4D] hover:text-[#4CAF4F] py-1 transition-colors"
              >
                {link.label}
              </a>
            ))}
          </nav>
          
          <div className="pt-4 border-t border-[#E8ECF2] flex flex-col gap-3">
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full text-center py-2.5 text-sm font-medium text-[#4CAF4F] border border-[#4CAF4F] rounded-md hover:bg-[#4CAF4F]/5 transition-colors"
            >
              Login
            </a>
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full text-center py-2.5 text-sm font-medium text-white bg-[#4CAF4F] hover:bg-[#388E3C] rounded-md shadow-sm transition-all"
            >
              Sign up / Get App
            </a>
          </div>
        </div>
      )}
    </header>
  );
}
