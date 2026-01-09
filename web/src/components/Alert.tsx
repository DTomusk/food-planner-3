type AlertType = 'info' | 'error' | 'success' | 'warning';

interface AlertProps {
    message: string;
    type?: AlertType;
}

const alertStyles: Record<AlertType, string> = {
    error: 'text-red-600 bg-red-50 border-red-200',
    success: 'text-green-600 bg-green-50 border-green-200',
    warning: 'text-yellow-600 bg-yellow-50 border-yellow-200',
    info: 'text-blue-600 bg-blue-50 border-blue-200'
};

const alertPrefixes: Record<AlertType, string> = {
    error: 'Error: ',
    success: 'Success: ',
    warning: 'Warning: ',
    info: 'Info: '
};

export default function Alert({ message, type = 'error' }: AlertProps) {
    const styleClass = alertStyles[type];
    const prefix = alertPrefixes[type];
    
    return (
        <div className={`font-bold p-3 rounded-md border ${styleClass}`}>
            {prefix}{message}
        </div>
    );
}