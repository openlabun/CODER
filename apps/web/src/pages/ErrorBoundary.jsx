import React from 'react';
export class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null, info: null };
  }
  static getDerivedStateFromError(error) {
    return { hasError: true };
  }
  componentDidCatch(error, info) {
    this.setState({ error, info });
    console.error("ErrorBoundary caught an error", error, info);
  }
  render() {
    if (this.state.hasError) {
      return <div style={{padding:'20px', color:'red'}}><h1>Something went wrong.</h1><pre>{this.state.error?.toString()}</pre><pre>{this.state.info?.componentStack}</pre></div>;
    }
    return this.props.children; 
  }
}
