import React from 'react';

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'glass' | 'interactive';
}

export const Card: React.FC<CardProps> = ({
  variant = 'default',
  children,
  className = '',
  ...props
}) => {
  const variantClass = variant === 'glass' ? 'ui-card-glass' : variant === 'interactive' ? 'ui-card-interactive' : '';
  return (
    <div className={`ui-card ${variantClass} ${className}`} {...props}>
      {children}
    </div>
  );
};
