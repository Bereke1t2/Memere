"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Send, CheckCircle2 } from "lucide-react";

export function Footer() {
  const [email, setEmail] = useState("");
  const [subscribed, setSubscribed] = useState(false);

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    if (email) {
      setSubscribed(true);
      setEmail("");
    }
  };

  return (
    <footer className="bg-[#263238] text-white pt-16 pb-14 text-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-12 lg:gap-8 pb-12 border-b border-[#3A454B]">
          
          {/* Left Column: Brand & Copyright & Socials (Col span 5) */}
          <div className="lg:col-span-5 flex flex-col justify-between space-y-8">
            
            {/* Logo */}
            <div className="space-y-4">
              <Link href="/" className="flex items-center gap-2.5 group">
                <div className="w-8 h-8 rounded-lg bg-[#4CAF4F] flex items-center justify-center text-white shadow-sm">
                  <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
                    <path d="M12 2L2 7l10 5 10-5-10-5zm0 9.2L4.5 7.4 12 3.6l7.5 3.8L12 11.2zm0 2.8L2 9v6l10 5 10-5V9l-10 5zm0 2.2l-7.5-3.8V11l7.5 3.8 7.5-3.8v2.6L12 16.2z"/>
                  </svg>
                </div>
                <div className="flex items-baseline gap-1">
                  <span className="text-2xl font-bold tracking-tight text-white">
                    Mirkuz
                  </span>
                  <span className="text-xs font-semibold text-[#4CAF4F]">
                    Nexcent
                  </span>
                </div>
              </Link>

              <div className="space-y-1 text-xs text-[#89939E]">
                <p>Copyright © 2026 Nexcent ltd. / Mirkuz.</p>
                <p>All rights reserved.</p>
              </div>
            </div>

            {/* Social Icons in dark rounded circles */}
            <div className="flex items-center gap-3">
              {/* Instagram */}
              <a
                href="https://instagram.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-full bg-[#3A454B] hover:bg-[#4CAF4F] flex items-center justify-center text-white transition-colors"
                aria-label="Instagram"
              >
                <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
                </svg>
              </a>

              {/* Dribbble */}
              <a
                href="https://dribbble.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-full bg-[#3A454B] hover:bg-[#4CAF4F] flex items-center justify-center text-white transition-colors"
                aria-label="Dribbble"
              >
                <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M12 0c-6.627 0-12 5.373-12 12s5.373 12 12 12 12-5.373 12-12-5.373-12-12-12zm10.169 11.235c-.179-.047-2.607-.662-5.267-.298.549 1.504 1.052 3.109 1.436 4.717 2.213-1.077 3.518-2.67 3.831-4.419zm-5.467-2.316c-1.393-2.427-2.923-4.577-3.151-4.891 3.528.948 6.273 3.655 7.159 7.234-1.282-.128-2.748-.063-4.008.207-1.127.243-2.201.611-3.208 1.084 1.157-1.312 2.23-2.522 3.208-3.634zm-7.618-4.631c.218.3 1.745 2.42 3.136 4.83-2.186.828-4.73 1.516-7.534 1.936 1.026-3.056 3.013-5.618 5.798-6.766zm-7.084 8.712c2.977-.442 5.679-1.171 7.994-2.05-.386 1.066-.807 2.155-1.267 3.256-4.512 1.34-8.825 1.385-9.284 1.385-.028-.276-.043-.556-.043-.841 0-.603.067-1.19.19-1.75zm1.884 4.545c.789-.009 4.316-.089 8.286-1.306 1.049 2.766 1.83 5.485 2.176 6.892-4.629.988-9.034-1.637-10.462-5.586zm12.029 4.606c-.328-1.332-1.071-3.929-2.073-6.58 2.404-.32 4.587.195 4.821.256-.475 2.784-2.222 5.09-4.748 6.324z"/>
                </svg>
              </a>

              {/* Twitter / X */}
              <a
                href="https://twitter.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-full bg-[#3A454B] hover:bg-[#4CAF4F] flex items-center justify-center text-white transition-colors"
                aria-label="Twitter"
              >
                <svg className="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
                  <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
                </svg>
              </a>

              {/* YouTube */}
              <a
                href="https://youtube.com"
                target="_blank"
                rel="noopener noreferrer"
                className="w-9 h-9 rounded-full bg-[#3A454B] hover:bg-[#4CAF4F] flex items-center justify-center text-white transition-colors"
                aria-label="YouTube"
              >
                <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24">
                  <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
                </svg>
              </a>
            </div>

          </div>

          {/* Column 2: Company (Col span 2) */}
          <div className="lg:col-span-2 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Company
            </h4>
            <ul className="space-y-2.5 text-xs text-[#ABBED1]">
              <li><a href="#hero" className="hover:text-white transition-colors">About us</a></li>
              <li><a href="#services" className="hover:text-white transition-colors">Blog</a></li>
              <li><a href="#testimonial" className="hover:text-white transition-colors">Contact us</a></li>
              <li><a href="#pricing" className="hover:text-white transition-colors">Pricing</a></li>
              <li><a href="#testimonial" className="hover:text-white transition-colors">Testimonials</a></li>
            </ul>
          </div>

          {/* Column 3: Support (Col span 2) */}
          <div className="lg:col-span-2 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Support
            </h4>
            <ul className="space-y-2.5 text-xs text-[#ABBED1]">
              <li><a href="#faq" className="hover:text-white transition-colors">Help center</a></li>
              <li><a href="#faq" className="hover:text-white transition-colors">Terms of service</a></li>
              <li><a href="#faq" className="hover:text-white transition-colors">Legal</a></li>
              <li><a href="#faq" className="hover:text-white transition-colors">Privacy policy</a></li>
              <li><a href="#faq" className="hover:text-white transition-colors">Status</a></li>
            </ul>
          </div>

          {/* Column 4: Stay up to date (Col span 3) */}
          <div className="lg:col-span-3 space-y-4">
            <h4 className="text-base font-bold text-white tracking-tight">
              Stay up to date
            </h4>
            
            <form onSubmit={handleSubscribe} className="relative">
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Your email address"
                required
                className="w-full bg-[#515B60]/70 text-white placeholder-[#ABBED1] text-xs px-4 py-3 rounded-md border border-transparent focus:border-[#4CAF4F] focus:outline-hidden pr-10 transition-colors"
              />
              <button
                type="submit"
                aria-label="Subscribe"
                className="absolute right-2.5 top-1/2 -translate-y-1/2 text-white hover:text-[#4CAF4F] transition-colors p-1 cursor-pointer"
              >
                <Send className="w-4 h-4" />
              </button>
            </form>

            {subscribed && (
              <div className="flex items-center gap-1.5 text-xs text-[#4CAF4F]">
                <CheckCircle2 className="w-3.5 h-3.5" />
                <span>Thank you for subscribing!</span>
              </div>
            )}
          </div>

        </div>

        {/* Bottom sub-footer */}
        <div className="pt-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-[#89939E]">
          <p>Designed with Figma Minimal Agency UI System • Built with Next.js & Tailwind CSS</p>
          <div className="flex items-center gap-4">
            <span>🇪🇹 Ethiopian Grade 12 National Curriculum</span>
          </div>
        </div>

      </div>
    </footer>
  );
}
