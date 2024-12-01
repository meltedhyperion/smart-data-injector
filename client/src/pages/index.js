import Image from 'next/image'
import idea from "../../public/idea.png"
import flow from "../../public/working-flow.png"
import highLevel from "../../public/high-level.png"
import Link from 'next/link'
function Home() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <h1>Smart Data Injector</h1>
      <Link href="https://drive.google.com/file/d/1p4KjEl6NGGZIQeFacH4SVF-s-1-SQqlK/view?usp=sharing">
        <button>Video Link</button>
      </Link> 
      <div>
      <Link href="/api-key-generator">
        <button>Generate API Key</button>
      </Link>

      <Link href="/data-injector">
        <button>Inject Data</button>
      </Link>
      </div>
      
      <br />
      <h2>The Idea</h2>
      <div className="p-8"><Image alt="profile" src={idea} height={800} width={1500} /></div>
      <br />
      <br />
      <h2>High level Idea</h2>
      <div className="p-8"><Image alt="profile" src={highLevel} height={800} width={700} /></div>
      <br />
      <h2>The Flow</h2>
      <div className="p-8"><Image alt="profile" src={flow} height={800} width={1500} /></div>
    </div>
  );
}
export default Home;
