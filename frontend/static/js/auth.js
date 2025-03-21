document.addEventListener('DOMContentLoaded', () => {
  // Tab switching functionality
  const tabs = document.querySelectorAll('.auth-tab');
  const forms = document.querySelectorAll('.auth-form');
  
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      const targetTab = tab.dataset.tab;
      
      // Update active tab
      tabs.forEach(t => t.classList.remove('active'));
      tab.classList.add('active');
      
      // Show corresponding form
      forms.forEach(form => {
        form.classList.remove('active');
        if (form.id === `${targetTab}Form`) {
          form.classList.add('active');
        }
      });
    });
  });
  
  // Toggle business fields visibility based on user type selection
  const userTypeRadios = document.querySelectorAll('input[name="userType"]');
  const businessFields = document.getElementById('businessFields');
  
  userTypeRadios.forEach(radio => {
    radio.addEventListener('change', () => {
      if (radio.value === 'business') {
        businessFields.style.display = 'block';
        // Make business fields required
        document.getElementById('company').required = true;
        document.getElementById('businessSize').required = true;
      } else {
        businessFields.style.display = 'none';
        // Make business fields not required
        document.getElementById('company').required = false;
        document.getElementById('businessSize').required = false;
      }
    });
  });
  
  // Login form submission
  const loginForm = document.getElementById('loginForm');
  const loginMessage = document.getElementById('loginMessage');
  
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Clear previous messages
    loginMessage.textContent = '';
    loginMessage.className = 'auth-message';
    
    // Show loading state
    const submitBtn = loginForm.querySelector('.submit-button');
    const btnText = submitBtn.querySelector('span');
    const originalText = btnText.textContent;
    btnText.textContent = 'Signing in...';
    submitBtn.disabled = true;
    
    try {
      const formData = {
        email: document.getElementById('loginEmail').value,
        password: document.getElementById('loginPassword').value
      };
      
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Login failed');
      }
      
      // Store auth token and user data
      localStorage.setItem('authToken', data.token);
      localStorage.setItem('userData', JSON.stringify(data.user));
      
      // Show success message
      loginMessage.textContent = 'Login successful! Redirecting...';
      loginMessage.classList.add('success');
      
      // Redirect to dashboard after successful login
      setTimeout(() => {
        window.location.href = '/dashboard';
      }, 1500);
      
    } catch (error) {
      console.error('Login error:', error);
      loginMessage.textContent = error.message || 'An error occurred during login';
      loginMessage.classList.add('error');
      
      // Reset button
      btnText.textContent = originalText;
      submitBtn.disabled = false;
    }
  });
  
  // Signup form submission
  const signupForm = document.getElementById('signupForm');
  const signupMessage = document.getElementById('signupMessage');
  
  signupForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Clear previous messages
    signupMessage.textContent = '';
    signupMessage.className = 'auth-message';
    
    // Show loading state
    const submitBtn = signupForm.querySelector('.submit-button');
    const btnText = submitBtn.querySelector('span');
    const originalText = btnText.textContent;
    btnText.textContent = 'Creating account...';
    submitBtn.disabled = true;
    
    try {
      const userType = document.querySelector('input[name="userType"]:checked').value;
      
      const formData = {
        email: document.getElementById('signupEmail').value,
        username: document.getElementById('signupUsername').value,
        password: document.getElementById('signupPassword').value,
        userType: userType,
        company: document.getElementById('company').value || '',
        businessSize: document.getElementById('businessSize').value || ''
      };
      
      const response = await fetch('/api/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Signup failed');
      }
      
      // Store auth token and user data
      localStorage.setItem('authToken', data.token);
      localStorage.setItem('userData', JSON.stringify(data.user));
      
      // Show success message
      signupMessage.textContent = 'Account created successfully! Redirecting...';
      signupMessage.classList.add('success');
      
      // Redirect to dashboard after successful signup
      setTimeout(() => {
        window.location.href = '/dashboard';
      }, 1500);
      
    } catch (error) {
      console.error('Signup error:', error);
      signupMessage.textContent = error.message || 'An error occurred during signup';
      signupMessage.classList.add('error');
      
      // Reset button
      btnText.textContent = originalText;
      submitBtn.disabled = false;
    }
  });
  
  // Check if user is already logged in
  const checkAuthStatus = () => {
    const token = localStorage.getItem('authToken');
    if (token) {
      // Redirect to home page if already logged in
      window.location.href = '/';
    }
  };
  
  // Check auth status on page load
  checkAuthStatus();
});
