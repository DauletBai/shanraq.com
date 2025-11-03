(function () {
  document.addEventListener("DOMContentLoaded", function () {
    var nav = document.querySelector('[data-nav-behavior="underline-toggle"]');
    if (!nav) {
      return;
    }

    var links = Array.prototype.slice.call(nav.querySelectorAll(".nav-link"));
    if (!links.length) {
      return;
    }

    var activeLink = null;
    var sectionEntries = [];
    var scrollOffset = 0;
    var ticking = false;

    var updateScrollOffset = function () {
      scrollOffset = (nav.offsetHeight || 0) + 40;
    };

    var refreshSectionOffsets = function () {
      sectionEntries = [];
      links.forEach(function (link) {
        var hash = link.hash;
        if (!hash || hash.length <= 1) {
          return;
        }
        var target = document.getElementById(hash.slice(1));
        if (!target) {
          return;
        }
        sectionEntries.push({
          link: link,
          target: target,
          top: target.getBoundingClientRect().top + window.pageYOffset,
        });
      });
      sectionEntries.sort(function (a, b) {
        return a.top - b.top;
      });
      resetIndicator();
    };

    var resetIndicator = function () {
      nav.style.setProperty("--nav-underline-width", "0px");
      nav.style.setProperty("--nav-underline-x", "0px");
      nav.style.setProperty("--nav-underline-opacity", "0");
      nav.style.setProperty("--nav-underline-color", "transparent");
    };

    var moveIndicator = function (link) {
      if (!link) {
        resetIndicator();
        return;
      }
      var width = link.offsetWidth;
      var left = link.offsetLeft;
      var color = window.getComputedStyle(link).color;
      nav.style.setProperty("--nav-underline-width", width + "px");
      nav.style.setProperty("--nav-underline-x", left + "px");
      nav.style.setProperty("--nav-underline-color", color);
      nav.style.setProperty("--nav-underline-opacity", "1");
    };

    var setActiveLink = function (link, options) {
      if (!link) {
        return;
      }
      var force = options && options.force;
      if (activeLink === link && !force) {
        moveIndicator(link);
        return;
      }
      activeLink = link;
      links.forEach(function (item) {
        item.classList.toggle("active", item === link);
      });
      nav.classList.add("nav-underline-enabled");
      moveIndicator(link);
    };

    var findLinkByHash = function (hash) {
      if (!hash) {
        return null;
      }
      for (var i = 0; i < links.length; i++) {
        if (links[i].hash === hash) {
          return links[i];
        }
      }
      return null;
    };

    var updateActiveFromScroll = function () {
      if (!sectionEntries.length) {
        return;
      }
      var scrollPosition = window.pageYOffset + scrollOffset;
      var candidate = sectionEntries[0];
      for (var i = 0; i < sectionEntries.length; i++) {
        if (scrollPosition >= sectionEntries[i].top) {
          candidate = sectionEntries[i];
        }
      }
      var nearBottom =
        window.innerHeight + window.pageYOffset >=
        document.documentElement.scrollHeight - 4;
      if (nearBottom) {
        candidate = sectionEntries[sectionEntries.length - 1];
      }
      if (candidate) {
        setActiveLink(candidate.link);
      }
    };

    var onScroll = function () {
      if (!ticking) {
        window.requestAnimationFrame(function () {
          updateActiveFromScroll();
          ticking = false;
        });
        ticking = true;
      }
    };

    var onResize = function () {
      updateScrollOffset();
      refreshSectionOffsets();
      window.requestAnimationFrame(function () {
        moveIndicator(activeLink);
        updateActiveFromScroll();
      });
    };

    nav.addEventListener("click", function (event) {
      var target = event.target.closest(".nav-link");
      if (!target || !nav.contains(target)) {
        return;
      }
      setActiveLink(target, { force: true });
    });

    window.addEventListener("hashchange", function () {
      var match = findLinkByHash(window.location.hash);
      if (match) {
        setActiveLink(match, { force: true });
      }
    });

    window.addEventListener("scroll", onScroll, { passive: true });
    window.addEventListener("resize", onResize);
    window.addEventListener("load", onResize);

    updateScrollOffset();
    refreshSectionOffsets();

    var initialLink = (function () {
      for (var i = 0; i < links.length; i++) {
        if (links[i].classList.contains("active")) {
          return links[i];
        }
      }
      return sectionEntries.length ? sectionEntries[0].link : links[0];
    })();

    if (initialLink) {
      setActiveLink(initialLink, { force: true });
    } else {
      resetIndicator();
    }

    updateActiveFromScroll();
  });
})();
