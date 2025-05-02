document.addEventListener('DOMContentLoaded', function() {

    // Mobile menu toggle
    const menuToggle = document.querySelector('.menu-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (menuToggle) {
        menuToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
            menuToggle.classList.toggle('active');
        });
    }

    // Close menu when clicking outside
    document.addEventListener('click', function(event) {
        const isClickInsideMenu = navMenu.contains(event.target);
        const isClickOnToggle = menuToggle.contains(event.target);
        
        if (!isClickInsideMenu && !isClickOnToggle && navMenu.classList.contains('active')) {
            navMenu.classList.remove('active');
            menuToggle.classList.remove('active');
        }
    });
    
    // Apply smooth scrolling to all anchor links (both in nav and CTA button)
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            // Skip if it's a different type of link
            if (this.getAttribute('href') === '#') return;
            
            smoothScrollToAnchor(this, e);
        });
    });

    // Navbar scroll effect
    const navbar = document.querySelector('.navbar');
    window.addEventListener('scroll', function() {
        if (window.scrollY > 50) {
            navbar.style.boxShadow = '0 5px 10px rgba(0, 0, 0, 0.1)';
            navbar.style.padding = '15px 0';
        } else {
            navbar.style.boxShadow = 'none';
            navbar.style.padding = '15px 0';
        }
    });

    // Initialize animations
    initAnimations();
});

// Smooth scroll function for anchor links
function smoothScrollToAnchor(anchorElement, event) {
    event.preventDefault();
    
    const targetId = anchorElement.getAttribute('href');
    const targetElement = document.querySelector(targetId);
    
    if (targetElement) {
        // Close mobile menu if open
        const navMenu = document.querySelector('.nav-menu');
        const menuToggle = document.querySelector('.menu-toggle');
        if (navMenu.classList.contains('active')) {
            navMenu.classList.remove('active');
            menuToggle.classList.remove('active');
        }

        // Calculate target position with offset for navbar
        const navbarHeight = document.querySelector('.navbar').offsetHeight;
        const targetPosition = targetElement.getBoundingClientRect().top + window.pageYOffset - navbarHeight;
        
        // Smooth scroll with easing function
        smoothScrollTo(targetPosition, 800);
    }
}

// Smooth scroll animation with easing
function smoothScrollTo(targetPosition, duration) {
    const startPosition = window.pageYOffset;
    const distance = targetPosition - startPosition;
    let startTime = null;
    
    function animation(currentTime) {
        if (startTime === null) startTime = currentTime;
        const timeElapsed = currentTime - startTime;
        const run = easeInOutQuad(timeElapsed, startPosition, distance, duration);
        window.scrollTo(0, run);
        if (timeElapsed < duration) requestAnimationFrame(animation);
    }
    
    // Easing function for smooth acceleration/deceleration
    function easeInOutQuad(t, b, c, d) {
        t /= d/2;
        if (t < 1) return c/2*t*t + b;
        t--;
        return -c/2 * (t*(t-2) - 1) + b;
    }
    
    requestAnimationFrame(animation);
}

function initAnimations() {
    const animatedElements = document.querySelectorAll('.feature-card, .download-card');

    checkIfInView();
    
    window.addEventListener('scroll', checkIfInView);
    
    function checkIfInView() {
        const windowHeight = window.innerHeight;
        const windowTopPosition = window.scrollY;
        const windowBottomPosition = windowTopPosition + windowHeight;
        
        animatedElements.forEach(element => {
            const elementHeight = element.offsetHeight;
            const elementTopPosition = element.offsetTop;
            const elementBottomPosition = elementTopPosition + elementHeight;
            
            if (
                (elementBottomPosition >= windowTopPosition) &&
                (elementTopPosition <= windowBottomPosition)
            ) {
                element.classList.add('animated');
            }
        });
    }
}
