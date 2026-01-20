import Link from "@/components/Link";

export default function NavBar() {
    return (
        <nav className="bg-white shadow p-4 mb-6">
            <div className="container mx-auto flex justify-between items-center">
                <h1 className="text-xl font-bold">FoodSmash</h1>
                <div>
                <Link onClick={() => {}} text="Home" />
                </div>
            </div>
        </nav>
    );
}