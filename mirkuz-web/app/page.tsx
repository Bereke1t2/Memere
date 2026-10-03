import React from "react";
import { Navbar } from "@/components/Navbar";
import { HeroSection } from "@/components/HeroSection";
import { ClientsBar } from "@/components/ClientsBar";
import { CommunitySection } from "@/components/CommunitySection";
import { FeatureSplitOne } from "@/components/FeatureSplitOne";
import { StatsCounter } from "@/components/StatsCounter";
import { FeatureSplitTwo } from "@/components/FeatureSplitTwo";
import { TestimonialSpotlight } from "@/components/TestimonialSpotlight";
import { BlogMarketingSection } from "@/components/BlogMarketingSection";
import { InteractiveExamDemo } from "@/components/InteractiveExamDemo";
import { StreamCurriculum } from "@/components/StreamCurriculum";
import { PricingSection } from "@/components/PricingSection";
import { FaqSection } from "@/components/FaqSection";
import { CtaBanner } from "@/components/CtaBanner";
import { Footer } from "@/components/Footer";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white text-[#4D4D4D] selection:bg-[#4CAF4F]/20 selection:text-[#263238] relative">
      {/* 1. Header / Navbar */}
      <Navbar />

      {/* Main Page Flow matching Figma Minimal Landing Page */}
      <main>
        {/* 2. Hero Section: Lessons and insights from 8 years */}
        <HeroSection />

        {/* 3. Our Clients: Trust & Institution Logos */}
        <ClientsBar />

        {/* 4. Manage your entire community in a single system: 3-Card Grid */}
        <CommunitySection />

        {/* 5. Feature Split 1: The unseen of spending three years at Pixelgrade */}
        <FeatureSplitOne />

        {/* 6. Stats & Achievements: Helping a local business reinvent itself */}
        <StatsCounter />

        {/* 7. Feature Split 2: How to design your site footer like we did */}
        <FeatureSplitTwo />

        {/* 8. Customer Spotlight: Tim Smith / Kalkidan Bekele Case Study */}
        <TestimonialSpotlight />

        {/* 9. Blog & Insights: Caring is the new marketing */}
        <BlogMarketingSection />

        {/* 10. Interactive Exam Simulator Micro-Experience */}
        <InteractiveExamDemo />

        {/* 11. Stream & Curriculum Breakdown */}
        <StreamCurriculum />

        {/* 12. Transparent ETB Pricing & Local Gateways */}
        <PricingSection />

        {/* 13. Frequently Asked Questions */}
        <FaqSection />

        {/* 14. Bottom High-Impact CTA: Pellentesque suscipit fringilla libero eu. */}
        <CtaBanner />
      </main>

      {/* 15. Minimal Dark Footer */}
      <Footer />
    </div>
  );
}
