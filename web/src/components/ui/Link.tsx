interface LinkProps {
    onClick: () => void;
    text: string;
}

export default function Link({ onClick, text }: LinkProps) {
    return <a 
    onClick={onClick}
    className='hover:underline cursor-pointer'
    >{text}</a>;
}